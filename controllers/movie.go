package controllers

import (
	"fmt"
	"io"
	"mediaserver/data"
	"mediaserver/services"
	"mediaserver/views"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const movieDir = "./movies"

func RegisterMovieEndpoints(e *gin.Engine) {
	e.GET("/movies", listMoviesHandler)
	e.GET("/stream/:id", streamMovieHandler)
	e.POST("/add", addMovieHandler)
	e.DELETE("/delete/:id", deleteMovieHandler)
	e.GET("/movies/count", getTotalMovieCount)
	e.GET("/movies/:id", getMovieById)
}

func EnsureMovieDirExists() error {
	return os.MkdirAll(movieDir, 0755)
}

func streamMovieHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Movie filename is required"})
		return
	}

	movie, err := services.GetMovieById(id, c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	filePath := path.Join(movieDir, movie.Filename)
	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve file info"})
		return
	}

	rangeHeader := c.GetHeader("Range")
	if rangeHeader != "" {
		rangeStart, rangeEnd := parseRange(rangeHeader, fileStat.Size())

		c.Header("Content-Type", "video/mp4")
		c.Header("Accept-Ranges", "bytes")
		c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", rangeStart, rangeEnd, fileStat.Size()))
		c.Status(http.StatusPartialContent)

		file.Seek(rangeStart, io.SeekStart)
		io.CopyN(c.Writer, file, rangeEnd-rangeStart+1)
	} else {
		c.Header("Content-Type", "video/mp4")
		c.Header("Content-Length", strconv.FormatInt(fileStat.Size(), 10))
		io.Copy(c.Writer, file)
	}
}

func listMoviesHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	page, err := strconv.Atoi(c.Query("page"))
	searchTerm := c.Query("search")
	if err != nil {
		page = 1
	}
	if page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		pageSize = 5
	}
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	var movies []*data.Movie
	searchTerm = strings.ToLower(searchTerm)
	db.Joins("LEFT JOIN omdb_movies ON omdb_movies.imdb_id = movies.omdb_id").
		Preload("OmdbMovie").
		Where("LOWER(omdb_movies.title) LIKE ?", "%"+searchTerm+"%").
		Offset(offset).Limit(pageSize).Find(&movies)
	var movieViews []*views.Movie
	for _, movie := range movies {
		movieViews = append(movieViews, movie.ToView())
	}

	c.JSON(http.StatusOK, gin.H{"movies": movieViews})
}

func getMovieById(c *gin.Context) {
	var id = c.Param("id")

	movie, err := services.GetMovieById(id, c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	c.JSON(http.StatusOK, movie.ToView())
}

func getTotalMovieCount(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var count int64
	db.Model(data.Movie{}).Count(&count)
	c.JSON(http.StatusOK, gin.H{"count": count})
}

func addMovieHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	imdbId := c.GetHeader("X-ImdbId")
	if imdbId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing ImdbId header"})
		return
	}

	filename := c.GetHeader("X-Filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing filename header"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file content"})
		return
	}

	var movie data.Movie
	movie.Filename = filename
	movie.Title = strings.Split(filename, ",")[0]
	movie.CreatedAt = time.Now()
	movie.ID = uuid.New()
	movie.OmdbId = imdbId

	omdbMovie, err := services.GetOmdbMovieById(imdbId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movie details from OMDB"})
		return
	}

	movie.OmdbMovie = omdbMovie

	tx := db.Begin()

	if err := tx.Create(&omdbMovie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create OMDB movie",
		})
		return
	}

	if err := tx.Create(&movie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create movie",
		})
		return
	}

	pathToSave := path.Join(movieDir, filename)
	err = os.WriteFile(pathToSave, fileBytes, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save movie"})
		tx.Rollback()
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to commit transaction",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Movie uploaded successfully"})
}

func deleteMovieHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id is required"})
		return
	}

	movie, err := services.GetMovieById(id, c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	pathToDelete := path.Join(movieDir, movie.Filename)

	tx := db.Begin()

	if err := tx.Delete(&movie).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie from database"})
		return
	}

	if _, err := os.Stat(pathToDelete); os.IsNotExist(err) {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	err = os.Remove(pathToDelete)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie deleted successfully"})
}

func parseRange(rangeHeader string, fileSize int64) (int64, int64) {
	var start, end int64
	fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)
	if end == 0 || end >= fileSize {
		end = fileSize - 1
	}
	return start, end
}
