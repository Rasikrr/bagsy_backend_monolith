// nolint: unused
package networks

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/network"
	"time"
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
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	CreatedBy   string  `json:"created_by"`
}

func toNetworkResponse(network *network.Network) *networkResponse {
	resp := &networkResponse{
		Code:        network.Code,
		Name:        network.Name,
		Description: network.Description,
		CreatedAt:   network.CreatedAt,
		CreatedBy:   network.CreatedBy,
	}

	if network.UpdatedAt != nil {
		updatedAt := network.UpdatedAt
		resp.UpdatedAt = updatedAt
	}

	return resp
}
