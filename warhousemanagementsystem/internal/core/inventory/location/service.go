package location

import (
	"context"
	"fmt"
	"strings"

	"github.com/mxV03/warhousemanagementsystem/ent"
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
