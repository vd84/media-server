package data

import (
	"mediaserver/views"

	"github.com/google/uuid"
)

type Series struct {
	ID        uuid.UUID `gorm:"primarykey"`
	OmdbMovie OmdbMovie `gorm:"foreignkey:OmdbId;references:IMDBID;constraint:OnDelete:CASCADE;"`
	OmdbId    string    `gorm:"foreignkey:IMDBID;references:IMDBID;constraint:OnDelete:CASCADE;"`
	Episodes  []Episode `gorm:"foreignkey:SeriesID;references:ID;constraint:OnDelete:CASCADE;"`
}

func (m *Series) ToView() *views.Series {
	return &views.Series{
		ID:        m.ID,
		OmdbMovie: m.OmdbMovie.ToView(),
	}
}
