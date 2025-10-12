package jobs

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/services/users"
	"github.com/Rasikrr/core/log"
)

type DeleteUnactivatedUsers struct {
	name           string
	schedule       string
	inactivePeriod time.Duration
	usersService   users.Service
}

func NewDeleteUnactivatedUsers(
	name string,
	schedule string,
	inactivePeriod time.Duration,
	usersService users.Service,
) *DeleteUnactivatedUsers {
	return &DeleteUnactivatedUsers{
		name:           name,
		schedule:       schedule,
		inactivePeriod: inactivePeriod,
		usersService:   usersService,
	}
}

func (j *DeleteUnactivatedUsers) Schedule() string {
	return j.schedule
}

func (j *DeleteUnactivatedUsers) Name() string {
	return j.name
}

func (j *DeleteUnactivatedUsers) Run() {
	var err error
	ctx := context.Background()
	l := log.With(log.String("name", j.Name()), log.String("schedule", j.Schedule()))

	defer func() {
		if err != nil {
			l.Error(ctx, "Failed to delete unactivated users", log.Err(err))
			return
		}
		l.Info(ctx, "Deleted unactivated users")
	}()
	l.Info(ctx, "Running job")
	err = j.usersService.DeleteUnactivatedUsers(ctx, j.inactivePeriod)
}
