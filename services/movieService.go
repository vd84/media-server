package services

import (
	"mediaserver/data"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetMovieById(id string, c *gin.Context) (*data.Movie, error) {
	var db = c.MustGet("db").(*gorm.DB)
	var movie data.Movie
	result := db.Preload("OmdbMovie").First(&movie, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &movie, nil
}
