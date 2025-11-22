package data

import (
	"mediaserver/views"
)

type Rating struct {
	Id      uint   `gorm:"unique;primaryKey;autoIncrement"`
	MovieID string `gorm:"foreignkey:IMDBID;references:IMDBID;constraint:OnDelete:CASCADE;"`
	Source  string
	Value   string
}

type OmdbMovie struct {
	Title      string
	Year       string
	Rated      string
	Released   string
	Runtime    string
	Genre      string
	Director   string
	Writer     string
	Actors     string
	Plot       string
	Language   string
	Country    string
	Awards     string
	Poster     string
	Ratings    []Rating `gorm:"foreignKey:MovieID;constraint:OnDelete:CASCADE;"`
	Metascore  string
	IMDBRating string
	IMDBVotes  string
	IMDBID     string `gorm:"primarykey"`
	Type       string
	DVD        string
	BoxOffice  string
	Production string
	Website    string
	Response   string
}

func (m *OmdbMovie) ToView() *views.OmdbMovie {
	vr := make([]views.Rating, 0, len(m.Ratings))
	for _, r := range m.Ratings {
		vr = append(vr, views.Rating{
			Source: r.Source,
			Value:  r.Value,
		})
	}

	return &views.OmdbMovie{
		Title:      m.Title,
		Year:       m.Year,
		Rated:      m.Rated,
		Released:   m.Released,
		Runtime:    m.Runtime,
		Genre:      m.Genre,
		Director:   m.Director,
		Writer:     m.Writer,
		Actors:     m.Actors,
		Plot:       m.Plot,
		Language:   m.Language,
		Country:    m.Country,
		Awards:     m.Awards,
		Poster:     m.Poster,
		Ratings:    vr,
		Metascore:  m.Metascore,
		IMDBRating: m.IMDBRating,
		IMDBVotes:  m.IMDBVotes,
		IMDBID:     m.IMDBID,
		Type:       m.Type,
		DVD:        m.DVD,
		BoxOffice:  m.BoxOffice,
		Production: m.Production,
		Website:    m.Website,
		Response:   m.Response,
	}
}
