package location

import (
	"fmt"
	"strings"
)

// ═══════════════════════════════════════════════════════════════
//              Coordinates
// ═══════════════════════════════════════════════════════════════

type Coordinates struct {
	Latitude  float64
	Longitude float64
}

func NewCoordinates(lat, lng float64) (Coordinates, error) {
	if lat < -90 || lat > 90 {
		return Coordinates{}, ErrInvalidLatitude
	}
	if lng < -180 || lng > 180 {
		return Coordinates{}, ErrInvalidLongitude
	}
	return Coordinates{Latitude: lat, Longitude: lng}, nil
}

func (c *Coordinates) IsZero() bool {
	return c.Latitude == 0 && c.Longitude == 0
}

func (c *Coordinates) String() string {
	return fmt.Sprintf("%f,%f", c.Latitude, c.Longitude)
}

// ═══════════════════════════════════════════════════════════════
//              Address
// ═══════════════════════════════════════════════════════════════

type Address struct {
	City     string
	Street   string
	Building string
	Details  string
}

func NewAddress(city, street, building, details string) (Address, error) {
	city = strings.TrimSpace(city)
	if city == "" {
		return Address{}, ErrCityRequired
	}

	return Address{
		City:     city,
		Street:   strings.TrimSpace(street),
		Building: strings.TrimSpace(building),
		Details:  strings.TrimSpace(details),
	}, nil
}

// Full — полный адрес одной строкой
func (a Address) Full() string {
	parts := []string{a.City}

	if a.Street != "" {
		streetPart := a.Street
		if a.Building != "" {
			streetPart += " " + a.Building
		}
		parts = append(parts, streetPart)
	}

	if a.Details != "" {
		parts = append(parts, a.Details)
	}

	return strings.Join(parts, ", ")
}

// Short — короткий адрес (без города)
func (a Address) Short() string {
	var parts []string

	if a.Street != "" {
		streetPart := a.Street
		if a.Building != "" {
			streetPart += " " + a.Building
		}
		parts = append(parts, streetPart)
	}

	if a.Details != "" {
		parts = append(parts, a.Details)
	}

	return strings.Join(parts, ", ")
}
