package domain

import (
	"time"

	"github.com/google/uuid"
)

type URLStatus int

const (
	StatusPending URLStatus = iota
	StatusSafe
	StatusMalicious
)

type URLCheck struct {
	ID        uuid.UUID `json:"id"`
	URL       string    `json:"url"`
	Status    URLStatus `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
