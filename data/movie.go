package data

import (
	"mediaserver/views"
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	ID        uuid.UUID `gorm:"primarykey"`
	Title     string
	Filename  string
	CreatedAt time.Time
	OmdbId    string
}

func (m *Movie) ToView() *views.Movie {
	return &views.Movie{
		ID:       m.ID,
		Title:    m.Title,
		Filename: m.Filename,
	}
}
