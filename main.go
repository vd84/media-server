package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

// Serve a static file (movie) based on its ID with optional streaming support.
func streamMovieHandler(w http.ResponseWriter, r *http.Request) {
	fileName := strings.TrimPrefix(r.URL.Path, "/stream/")
	if fileName == "" {
		http.Error(w, "Movie ID is required", http.StatusBadRequest)
		return
	}

	filePath := path.Join("./movies", fileName)
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		http.Error(w, "Unable to retrieve file info", http.StatusInternalServerError)
		return
	}

	rangeHeader := r.Header.Get("Range")
	if rangeHeader != "" {
		rangeStart, rangeEnd := parseRange(rangeHeader, fileStat.Size())

		w.Header().Set("Content-Type", "video/mp4")
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set(
			"Content-Range",
			fmt.Sprintf("bytes %d-%d/%d", rangeStart, rangeEnd, fileStat.Size()),
		)
		w.WriteHeader(http.StatusPartialContent)

		file.Seek(rangeStart, io.SeekStart)
		io.CopyN(w, file, rangeEnd-rangeStart+1)
	} else {
		w.Header().Set("Content-Type", "video/mp4")
		w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))
		io.Copy(w, file)
	}
}

func parseRange(rangeHeader string, fileSize int64) (int64, int64) {
	var start, end int64
	fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)
	if end == 0 || end >= fileSize {
		end = fileSize - 1
	}
	return start, end
}

func listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	// Path to the movies directory
	movieDir := "./movies"

	// Read directory contents
	files, err := os.ReadDir(movieDir)
	if err != nil {
		http.Error(w, "Unable to list movies", http.StatusInternalServerError)
		return
	}

	var movies []string
	for _, file := range files {
		if !file.IsDir() { // Only include files
			movies = append(movies, file.Name())
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"movies": movies,
	})
}

func addMovieHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.Header.Get("X-Filename")
	if filename == "" {
		http.Error(w, "Missing filename", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the file content into memory
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file content", http.StatusInternalServerError)
		return
	}

	pathToSaveTo := "./movies/" + filename
	err = os.WriteFile(pathToSaveTo, fileBytes, 0644)
	if err != nil {
		http.Error(w, "Failed to save movie", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Movie uploaded successfully"))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Filename")

		// Handle preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register your handlers
	mux.HandleFunc("/stream/", streamMovieHandler)
	mux.HandleFunc("/movies", listMoviesHandler)
	mux.HandleFunc("/add", addMovieHandler)

	// Wrap the ServeMux with the CORS middleware
	corsHandler := corsMiddleware(mux)

	fmt.Println("Starting media server on http://localhost:8080...")
	err := http.ListenAndServe(":8080", corsHandler)
	if err != nil {
		return
	}
}
