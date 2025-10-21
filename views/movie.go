package views

import (
	"github.com/google/uuid"
)

type Movie struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Filename string    `json:"filename"`
}
