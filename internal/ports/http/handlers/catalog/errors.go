package catalog

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
	domainCatalog "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/catalog"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
)

var catalogErrors = httputil.ErrorMap{
	billing.ErrSubscriptionSuspended: {Code: http.StatusForbidden, Message: "subscription_suspended"},
	identity.ErrPermissionDenied:     {Code: http.StatusForbidden, Message: "permission_denied"},

	location.ErrLocationNotFound: {Code: http.StatusNotFound, Message: "location_not_found"},
	location.ErrLocationInactive: {Code: http.StatusUnprocessableEntity, Message: "location_inactive"},

	domainCatalog.ErrServiceNotFound:          {Code: http.StatusNotFound, Message: "service_not_found"},
	domainCatalog.ErrServiceDeleted:           {Code: http.StatusGone, Message: "service_deleted"},
	domainCatalog.ErrServiceInactive:          {Code: http.StatusUnprocessableEntity, Message: "service_inactive"},
	domainCatalog.ErrServiceNameRequired:      {Code: http.StatusBadRequest, Message: "service_name_required"},
	domainCatalog.ErrServiceInvalidColor:      {Code: http.StatusBadRequest, Message: "service_invalid_color"},
	domainCatalog.ErrServiceCategoryNotFound:  {Code: http.StatusNotFound, Message: "service_category_not_found"},
	domainCatalog.ErrCategoryMismatch:         {Code: http.StatusBadRequest, Message: "category_mismatch"},
	domainCatalog.ErrEmployeeLocationMismatch: {Code: http.StatusUnprocessableEntity, Message: "employee_location_mismatch"},

	identity.ErrEmployeeNotFound:    {Code: http.StatusNotFound, Message: "employee_not_found"},
	identity.ErrEmployeeCannotServe: {Code: http.StatusUnprocessableEntity, Message: "employee_cannot_serve"},

	shared.ErrInvalidDuration: {Code: http.StatusBadRequest, Message: "invalid_duration"},
	shared.ErrNegativeAmount:  {Code: http.StatusBadRequest, Message: "negative_amount"},
	shared.ErrInvalidMoney:    {Code: http.StatusBadRequest, Message: "invalid_money"},
}
