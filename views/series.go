package views

import "github.com/google/uuid"

type Series struct {
	ID        uuid.UUID  `json:"id"`
	OmdbMovie *OmdbMovie `json:"omdbMovie,omitempty"`
}
