package point

import domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"

var (
	ErrPointNotFound                = domainErr.NewNotFoundError("point(s) not found", nil)
	ErrPointCategoryNotFound        = domainErr.NewNotFoundError("point category(ies) not found", nil)
	ErrPointAlreadyExists           = domainErr.NewInvalidInputError("point with same code already exists", nil)
	ErrPointCategoryServiceNotFound = domainErr.NewNotFoundError("point category(ies) services not found", nil)
)
