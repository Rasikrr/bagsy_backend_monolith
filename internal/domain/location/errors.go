package location

import "errors"

var (
	ErrLocationNotFound     = errors.New("location not found")
	ErrLocationDeleted      = errors.New("location deleted")
	ErrLocationInactive     = errors.New("location is inactive")
	ErrNameRequired         = errors.New("name is required")
	ErrCityRequired         = errors.New("city required")
	ErrInvalidScheduleType  = errors.New("invalid schedule type for location")
	ErrScheduleTypeRequired = errors.New("schedule type is required")
	ErrInvalidLatitude      = errors.New("latitude must be between -90 and 90")
	ErrInvalidLongitude     = errors.New("longitude must be between -180 and 180")
)

var (
	ErrCategoryNotFound     = errors.New("location_category not found")
	ErrCategoryNameRequired = errors.New("location_category name is required")
)
