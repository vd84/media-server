package views

import "github.com/google/uuid"

type Season struct {
	ID       uuid.UUID `json:"id"`
	Index    int       `json:"index"`
	Title    string    `json:"title"`
	Series   *Series   `json:"series,omitempty"`
	Episodes []Episode `json:"episodes,omitempty"`
}
