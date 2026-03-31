package access

import "github.com/google/uuid"

type OrganizationInfo struct {
	ID     uuid.UUID
	Active bool
	Name   string
}
