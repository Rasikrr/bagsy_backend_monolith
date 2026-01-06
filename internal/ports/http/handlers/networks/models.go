package networks

import "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/command"

//go:generate easyjson -all models.go

type createNetworkRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *createNetworkRequest) toDomain() *command.CreateNetworkCommand {
	return &command.CreateNetworkCommand{
		Name:        c.Name,
		Description: c.Description,
	}
}
