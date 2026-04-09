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
	ID uuid.UUID
	URL string
	Status URLStatus
	CreatedAt time.Time
}



