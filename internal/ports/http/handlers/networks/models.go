// nolint: unused
package networks

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/network"
)

//go:generate easyjson -all models.go

type createNetworkRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *createNetworkRequest) toDomain() *network.CreateNetworkCommand {
	return &network.CreateNetworkCommand{
		Name:        c.Name,
		Description: c.Description,
	}
}

type networkResponse struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   *string `json:"updated_at,omitempty"`
	CreatedBy   string  `json:"created_by"`
}

func toNetworkResponse(network *network.Network) *networkResponse {
	resp := &networkResponse{
		Code:        network.Code,
		Name:        network.Name,
		Description: network.Description,
		CreatedAt:   network.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		CreatedBy:   network.CreatedBy,
	}

	if network.UpdatedAt != nil {
		updatedAt := network.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
		resp.UpdatedAt = &updatedAt
	}

	return resp
}
