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
	if err := applyNamePatch(loc, input); err != nil {
		return err
	}
	if err := applyPhonePatch(loc, input); err != nil {
		return err
	}
	if err := applyLocationAddressPatch(loc, input); err != nil {
		return err
	}
	if err := applyActivePatch(loc, input); err != nil {
		return err
	}
	if err := applyScheduleTypePatch(loc, orgCtx, u, input); err != nil {
		return err
	}
	return applySlotDurationPatch(loc, input)
}

func applyNamePatch(loc *location.Location, input UpdateLocationInput) error {
	if input.Name == nil {
		return nil
	}
	return loc.UpdateInfo(*input.Name, loc.Description)
}

func applyPhonePatch(loc *location.Location, input UpdateLocationInput) error {
	if input.Phone == nil {
		return nil
	}
	phone, err := shared.NewPhone(*input.Phone)
	if err != nil {
		return err
	}
	return loc.ChangePhone(&phone)
}

func applyLocationAddressPatch(loc *location.Location, input UpdateLocationInput) error {
	if input.Address == nil && input.Latitude == nil && input.Longitude == nil {
		return nil
	}
	addr, coords, err := resolveAddressPatch(loc, input)
	if err != nil {
		return err
	}
	return loc.SetAddress(addr, coords)
}

func applyActivePatch(loc *location.Location, input UpdateLocationInput) error {
	if input.Active == nil {
		return nil
	}
	if *input.Active {
		return loc.Activate()
	}
	return loc.Deactivate()
}

func applyScheduleTypePatch(loc *location.Location, orgCtx *access.OrgContext, u *UseCase, input UpdateLocationInput) error {
	if input.ScheduleType == nil {
		return nil
	}
	st, err := u.resolveScheduleType(orgCtx.Plan.Code, *input.ScheduleType)
	if err != nil {
		return err
	}
	return loc.ChangeScheduleType(st)
}

func applySlotDurationPatch(loc *location.Location, input UpdateLocationInput) error {
	if input.SlotDurationMinutes == nil {
		return nil
	}
	duration, err := shared.NewDuration(*input.SlotDurationMinutes)
	if err != nil {
		return err
	}
	return loc.ChangeSlotDuration(duration)
}

func resolveAddressPatch(loc *location.Location, input UpdateLocationInput) (*location.Address, *location.Coordinates, error) {
	addr, err := resolveAddress(loc, input)
	if err != nil {
		return nil, nil, err
	}
	coords, err := resolveCoordinates(loc, input)
	if err != nil {
		return nil, nil, err
	}
	return addr, coords, nil
}

func resolveAddress(loc *location.Location, input UpdateLocationInput) (*location.Address, error) {
	if input.Address == nil {
		return loc.Address, nil
	}
	city, street, building, details := currentAddressFields(loc)
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
		return nil, err
	}
	return &a, nil
}

func currentAddressFields(loc *location.Location) (city, street, building, details string) {
	if loc.Address == nil {
		return "", "", "", ""
	}
	return loc.Address.City, loc.Address.Street, loc.Address.Building, loc.Address.Details
}

func resolveCoordinates(loc *location.Location, input UpdateLocationInput) (*location.Coordinates, error) {
	if input.Latitude == nil && input.Longitude == nil {
		return loc.Coordinates, nil
	}
	lat, lng := float64(0), float64(0)
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
		return nil, err
	}
	return &c, nil
}
