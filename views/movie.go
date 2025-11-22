package views

import (
	"github.com/google/uuid"
)

type Movie struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Filename string    `json:"filename"`
	Year     string    `json:"year,omitempty"`
	IMDBID   string    `json:"imdbID,omitempty"`
	Type     string    `json:"type,omitempty"`
	Poster   string    `json:"poster,omitempty"`
}
