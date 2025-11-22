package data

import (
	"mediaserver/views"

	"github.com/google/uuid"
)

type Episode struct {
	ID       uuid.UUID `gorm:"primarykey"`
	Index    int
	Series   Series    `gorm:"foreignkey:SeriesID;references:ID;constraint:OnDelete:CASCADE;"`
	SeriesID uuid.UUID `gorm:"foreignkey:SeriesID;references:ID;constraint:OnDelete:CASCADE;"`
}

func (m *Episode) ToView() *views.Episode {
	return &views.Episode{
		ID:    m.ID,
		Index: m.Index,
	}
}
