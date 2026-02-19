package organization

import "errors"

var (
	ErrOrganizationNotFound          = errors.New("organization not found")
	ErrOrganizationProfileIncomplete = errors.New("organization profile incomplete")
	ErrOrganizationNameRequired      = errors.New("organization name is required")
	ErrEmptySlug                     = errors.New("empty slug of organization")
	ErrOrganizationDeleted           = errors.New("organization deleted")
	ErrOrganizationInactive          = errors.New("organization is inactive")
)
