package service

import domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"

var (
	ErrServiceNotFound            = domainErr.NewNotFoundError("service(s) not found", nil)
	ErrServiceCategoryNotFound    = domainErr.NewNotFoundError("service category(ies) not found", nil)
	ErrServiceSubcategoryNotFound = domainErr.NewNotFoundError("service subcategory(ies) not found", nil)
)
