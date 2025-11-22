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
	OmdbMovie OmdbMovie     `gorm:"foreignkey:OmdbId;references:IMDBID;constraint:OnDelete:CASCADE;nullable"`
	OmdbId    string        `gorm:"foreignkey:IMDBID;references:IMDBID;constraint:OnDelete:CASCADE;nullable"`
	Episode   Episode       `gorm:"foreignkey:EpisodeID;references:ID;constraint:OnDelete:CASCADE;nullable"`
	EpisodeID uuid.NullUUID `gorm:"foreignkey:ID;references:ID;constraint:OnDelete:CASCADE;nullable"`
}

func (m *Movie) ToView() *views.Movie {
	return &views.Movie{
		ID:        m.ID,
		Title:     m.Title,
		Filename:  m.Filename,
		OmdbMovie: m.OmdbMovie.ToView(),
	}
}
