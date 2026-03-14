package location

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/access"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
)

func (u *UseCase) UpdateLocation(ctx context.Context, orgCtx *access.OrgContext, input UpdateLocationInput) error {
	loc, err := u.locationRepo.GetByID(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("get location: %w", err)
	}

	if err = u.policy.CanManageLocation(orgCtx, loc); err != nil {
		return err
	}

	if err = applyLocationPatch(loc, orgCtx, u, input); err != nil {
		return err
	}

	if err = u.locationRepo.Save(ctx, loc); err != nil {
		return fmt.Errorf("save location: %w", err)
	}
	return nil
}

func applyLocationPatch(loc *location.Location, orgCtx *access.OrgContext, u *UseCase, input UpdateLocationInput) error {
	if input.Name != nil {
		if err := loc.UpdateInfo(*input.Name, loc.Description); err != nil {
			return err
		}
	}

	if input.Phone != nil {
		phone, err := shared.NewPhone(*input.Phone)
		if err != nil {
			return err
		}
		if err = loc.ChangePhone(&phone); err != nil {
			return err
		}
	}

	if input.Address != nil || input.Latitude != nil || input.Longitude != nil {
		addr, coords, err := resolveAddressPatch(loc, input)
		if err != nil {
			return err
		}
		if err = loc.SetAddress(addr, coords); err != nil {
			return err
		}
	}

	if input.Active != nil {
		if *input.Active {
			if err := loc.Activate(); err != nil {
				return err
			}
		} else {
			if err := loc.Deactivate(); err != nil {
				return err
			}
		}
	}

	if input.ScheduleType != nil {
		st, err := u.resolveScheduleType(orgCtx.Plan.Code, *input.ScheduleType)
		if err != nil {
			return err
		}
		if err = loc.ChangeScheduleType(st); err != nil {
			return err
		}
	}

	if input.SlotDurationMinutes != nil {
		duration, err := shared.NewDuration(*input.SlotDurationMinutes)
		if err != nil {
			return err
		}
		if err = loc.ChangeSlotDuration(duration); err != nil {
			return err
		}
	}

	return nil
}

func resolveAddressPatch(loc *location.Location, input UpdateLocationInput) (*location.Address, *location.Coordinates, error) {
	var addr *location.Address
	if input.Address != nil {
		// merge with current address values for nil fields
		city := ""
		street := ""
		building := ""
		details := ""

		if loc.Address != nil {
			city = loc.Address.City
			street = loc.Address.Street
			building = loc.Address.Building
			details = loc.Address.Details
		}

		if input.Address.City != nil {
			city = *input.Address.City
		}
		if input.Address.Street != nil {
			street = *input.Address.Street
		}
		if input.Address.Building != nil {
			building = *input.Address.Building
		}
		if input.Address.Details != nil {
			details = *input.Address.Details
		}

		a, err := location.NewAddress(city, street, building, details)
		if err != nil {
			return nil, nil, err
		}
		addr = &a
	} else {
		addr = loc.Address
	}

	var coords *location.Coordinates
	if input.Latitude != nil || input.Longitude != nil {
		lat := float64(0)
		lng := float64(0)
		if loc.Coordinates != nil {
			lat = loc.Coordinates.Latitude
			lng = loc.Coordinates.Longitude
		}
		if input.Latitude != nil {
			lat = *input.Latitude
		}
		if input.Longitude != nil {
			lng = *input.Longitude
		}
		c, err := location.NewCoordinates(lat, lng)
		if err != nil {
			return nil, nil, err
		}
		coords = &c
	} else {
		coords = loc.Coordinates
	}

	return addr, coords, nil
}
