package data

import (
	"mediaserver/views"

	"github.com/google/uuid"
)

type Season struct {
	ID       uuid.UUID `gorm:"primarykey"`
	Index    int
	Title    string
	Episodes []Episode `gorm:"foreignkey:EpisodeID;references:ID;constraint:OnDelete:CASCADE;nullable"`
	Series   Series    `gorm:"foreignkey:SeriesID;references:ID;constraint:OnDelete:CASCADE;"`
	SeriesID uuid.UUID `gorm:"foreignkey:SeriesID;references:ID;constraint:OnDelete:CASCADE;"`
}

func (m *Season) ToView() *views.Season {
	return &views.Season{
		ID:     m.ID,
		Index:  m.Index,
		Title:  m.Title,
		Series: m.Series.ToView(),
		Episodes: func() []views.Episode {
			episodes := make([]views.Episode, len(m.Episodes))
			for i, ep := range m.Episodes {
				episodes[i] = *ep.ToView()
			}
			return episodes
		}(),
	}
}
