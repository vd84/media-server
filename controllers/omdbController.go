package controllers

import (
	"mediaserver/services"

	"github.com/gin-gonic/gin"
)

func RegiserOmdbEndpoints(e *gin.Engine) {
	e.GET("/omdb/search", searchOmdbMovieHandler)
	e.GET("/omdb/movie/:id", getOmdbMovieByIdHandler)
}

func searchOmdbMovieHandler(c *gin.Context) {
	searchTerm := c.Query("searchTerm")

	res, err := services.GetMoviesBySearch(searchTerm)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch data from OMDB"})
		return
	}

	c.JSON(200, res)
}

func getOmdbMovieByIdHandler(c *gin.Context) {
	movieId := c.Param("id")

	res, err := services.GetMovieById(movieId)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch data from OMDB"})
		return
	}

	c.JSON(200, res)
}
