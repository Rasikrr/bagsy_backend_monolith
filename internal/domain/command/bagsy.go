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

//
//// UpdateBagsyCommand - команда для обновления брони
//type UpdateBagsyCommand struct {
//	// ID брони
//	ID uuid.UUID `json:"id" validate:"required"`
//
//	// StartAt - новое время начала (опционально)
//	StartAt *time.Time `json:"start_at,omitempty"`
//
//	// MasterPhone - перенос на другого мастера (опционально)
//	// При переносе цена пересчитывается из нового master_services
//	MasterPhone *string `json:"master_phone,omitempty" validate:"omitempty,min=10,max=15"`
//
//	// Comment - обновление комментария (опционально)
//	Comment *string `json:"comment,omitempty" validate:"omitempty,max=500"`
//}
//
//// CancelBagsyCommand - команда для отмены брони
//type CancelBagsyCommand struct {
//	ID     uuid.UUID `json:"id" validate:"required"`
//	Reason string    `json:"reason,omitempty" validate:"max=500"` // Будет записан в reject_reason
//}
