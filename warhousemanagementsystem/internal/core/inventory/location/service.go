package location

import (
	"context"
	"fmt"
	"strings"

	"github.com/mxV03/warhousemanagementsystem/ent"
	"github.com/mxV03/warhousemanagementsystem/ent/location"
)

var (
	ErrLocationNotFound = fmt.Errorf("location not found")
	ErrInvalidCode      = fmt.Errorf("invalid location code")
	ErrInvalidName      = fmt.Errorf("invalid location name")
	ErrLocationExists   = fmt.Errorf("location already exists")
)

type LocationService struct {
	client *ent.Client
}

func NewLocationService(client *ent.Client) *LocationService {
	return &LocationService{client: client}
}

type LocationDTO struct {
	ID   int
	Code string
	Name string
}

func (s *LocationService) CreateLocation(ctx context.Context, code, name string) (*LocationDTO, error) {
	code = strings.TrimSpace(code)
	name = strings.TrimSpace(name)
	if code == "" {
		return nil, ErrInvalidCode
	}
	if name == "" {
		return nil, ErrInvalidName
	}

	exists, err := s.client.Location.Query().Where(location.Code(code)).Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("checking location existence: %w", err)
	}
	if exists {
		return nil, ErrLocationExists
	}
	loc, err := s.client.Location.Create().
		SetCode(code).
		SetName(name).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating location: %w", err)
	}
	return &LocationDTO{
		ID:   loc.ID,
		Code: loc.Code,
		Name: loc.Name,
	}, nil
}

func (s *LocationService) GetLocationByCode(ctx context.Context, code string) (*LocationDTO, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, ErrInvalidCode
	}

	loc, err := s.client.Location.Query().Where(location.Code(code)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrLocationNotFound
		}
		return nil, fmt.Errorf("retrieving location: %w", err)
	}

	return &LocationDTO{
		ID:   loc.ID,
		Code: loc.Code,
		Name: loc.Name,
	}, nil
}

func (s *LocationService) ListLocations(ctx context.Context, limit int) ([]*LocationDTO, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	locations, err := s.client.Location.Query().Order(ent.Asc(location.FieldCode)).Limit(limit).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing locations: %w", err)
	}
	out := make([]*LocationDTO, 0, len(locations))
	for _, loc := range locations {
		out = append(out, &LocationDTO{
			ID:   loc.ID,
			Code: loc.Code,
			Name: loc.Name,
		})
	}
	return out, nil
}

func (s *LocationService) DeleteLocationByCode(ctx context.Context, code string) error {
	code = strings.TrimSpace(code)
	if code == "" {
		return ErrInvalidCode
	}

	deleted, err := s.client.Location.Delete().Where(location.Code(code)).Exec(ctx)
	if err != nil {
		return fmt.Errorf("deleting location by code: %w", err)
	}

	if deleted == 0 {
		return ErrLocationNotFound
	}
	return nil
}
