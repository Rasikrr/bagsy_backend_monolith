package command

import (
	"time"

	"github.com/google/uuid"
)

// CreateBagsyCommand - команда для создания брони
type CreateBagsyCommand struct {
	ServiceID   uuid.UUID
	MasterPhone string

	StartAt time.Time

	ClientPhone string
	Name        string
	Surname     string
	Comment     *string
}
