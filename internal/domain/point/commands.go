package point

import "github.com/google/uuid"

type CreatePointCommand struct {
	Name        string
	Description *string
	CategoryID  int
	Address     Address
	Schedule    Schedule
	PhotoIDs    []uuid.UUID
}
