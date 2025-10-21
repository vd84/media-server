package data

import (
	"github.com/google/uuid"
	"mediaserver/views"
	"time"
)

type Movie struct {
	ID        uuid.UUID `gorm:"primarykey"`
	Title     string
	Filename  string
	CreatedAt time.Time
}

func (m *Movie) ToView() *views.Movie {
	return &views.Movie{
		ID:       m.ID,
		Title:    m.Title,
		Filename: m.Filename,
	}
}
