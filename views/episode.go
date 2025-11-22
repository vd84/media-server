package views

import "github.com/google/uuid"

type Episode struct {
	ID        uuid.UUID `json:"id"`
	Index     int       `json:"index"`
	OmdbMovie OmdbMovie `json:"omdbMovie,omitempty"`
}
