package users

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/query"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
)

//go:generate easyjson -all models.go

type getUsersRequest struct {
	PointCode   *string  `query:"point_code" validate:""`
	NetworkCode *string  `query:"network_code"`
	Roles       []string `query:"roles"`
	Phones      []string `query:"phones"`
	Limit       uint64   `query:"limit" validate:"max=100"`
	Offset      uint64   `query:"offset" validate:"min=0"`
	SortBy      string   `query:"sort_by" validate:"oneof=phone name surname point_code network_code created_at updated_at"`
	SortOrder   string   `query:"sort_order" validate:"oneof=asc desc"`
}

func (r *getUsersRequest) GetQueryParameters(req *http.Request) error {
	q := req.URL.Query()

	if pointCode := q.Get("point_code"); pointCode != "" {
		r.PointCode = &pointCode
	}

	if networkCode := q.Get("network_code"); networkCode != "" {
		r.NetworkCode = &networkCode
	}

	if len(q.Get("roles")) > 0 {
		r.Roles = strings.Split(q.Get("roles"), ",")
	}

	if phones := q["phones"]; len(phones) > 0 {
		r.Phones = phones
	}

	r.Limit = 20
	if limitStr := q.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			r.Limit = uint64(limit)
		}
	}

	r.Offset = 0
	if offsetStr := q.Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset > 0 {
			r.Offset = uint64(offset)
		}
	}

	r.SortBy = "created_at"
	if sortBy := q.Get("sort_by"); sortBy != "" {
		r.SortBy = sortBy
	}

	r.SortOrder = "asc"
	if sortOrder := q.Get("sort_order"); sortOrder != "" {
		r.SortOrder = sortOrder
	}

	return nil
}

func (r *getUsersRequest) Validate() error {
	err := request.GetValidator().Struct(r)
	if err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

func (r *getUsersRequest) toFilter() (*query.UserFilter, error) {
	q := query.UserFilter{
		PointCode:   r.PointCode,
		NetworkCode: r.NetworkCode,
		Phones:      r.Phones,
		Limit:       r.Limit,
		Offset:      r.Offset,
		SortBy:      r.SortBy,
	}
	sortOrder, err := enum.SortOrderString(r.SortOrder)
	if err != nil {
		return nil, domainErr.NewValidationError("sort_order contains an invalid value").
			WithDetail("sort_order", r.SortOrder)
	}
	q.SortOrder = sortOrder

	roles := make([]enum.Role, len(r.Roles))
	for i, role := range r.Roles {
		roleEnum, err := enum.RoleString(role)
		if err != nil {
			return nil, domainErr.NewValidationError("roles contains invalid item").
				WithDetail("role", role)
		}
		roles[i] = roleEnum
	}
	q.Roles = roles

	return &q, nil
}

type userDTO struct {
	Phone       string  `json:"phone"`
	Role        string  `json:"role"`
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	PointCode   *string `json:"point_code,omitempty"`
	NetworkCode *string `json:"network_code,omitempty"`
	Active      bool    `json:"active"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   *string `json:"updated_at,omitempty"`
}

func toUserDTO(user *entity.User) userDTO {
	dto := userDTO{
		Phone:       user.Phone,
		Role:        user.Role.String(),
		Name:        user.Name,
		Surname:     user.Surname,
		PointCode:   user.PointCode,
		NetworkCode: user.NetworkCode,
		Active:      user.Active,
		CreatedAt:   user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if user.UpdatedAt != nil {
		updatedAt := user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
		dto.UpdatedAt = &updatedAt
	}

	return dto
}

type getUsersResponse struct {
	Users []userDTO `json:"users"`
	Count int       `json:"count"`
}

func toGetUsersResponse(users []*entity.User) getUsersResponse {
	dtos := make([]userDTO, 0, len(users))
	for _, user := range users {
		dtos = append(dtos, toUserDTO(user))
	}

	return getUsersResponse{
		Users: dtos,
		Count: len(dtos),
	}
}
