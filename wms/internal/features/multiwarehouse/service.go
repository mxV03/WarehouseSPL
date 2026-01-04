//go:build multiwarehouse

package multiwarehouse

import (
	"context"
	"fmt"
	"strings"

	"github.com/mxV03/wms/ent"
	"github.com/mxV03/wms/ent/location"
	"github.com/mxV03/wms/ent/warehouse"
	"github.com/mxV03/wms/ent/warehouselocation"
)

var (
	ErrInvalidWarehouseCode    = fmt.Errorf("invalid warehouse code")
	ErrInvalidLocationCode     = fmt.Errorf("invalid location code")
	ErrWarehouseExists         = fmt.Errorf("warehouse already exists")
	ErrWarehouseNotFound       = fmt.Errorf("warehouse not found")
	ErrLocationNotFound        = fmt.Errorf("location not found")
	ErrLocationAlreadyAssigned = fmt.Errorf("location already assigned to a warehouse")
)

type MultiwarehouseService struct {
	client *ent.Client
}

func NewMultiwarehouseService(client *ent.Client) *MultiwarehouseService {
	return &MultiwarehouseService{
		client: client,
	}
}

func (s *MultiwarehouseService) CreateWarehouse(ctx context.Context, code, name string) (*ent.Warehouse, error) {
	code = strings.TrimSpace(code)
	name = strings.TrimSpace(name)
	if code == "" {
		return nil, ErrInvalidWarehouseCode
	}

	exists, err := s.client.Warehouse.Query().
		Where(warehouse.Code(code)).
		Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("check warehouse existence: %w", err)
	}
	if exists {
		return nil, ErrWarehouseExists
	}

	w, err := s.client.Warehouse.Create().
		SetCode(code).
		SetName(name).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, ErrWarehouseExists
		}
		return nil, fmt.Errorf("create warehouse: %w", err)
	}
	return w, nil
}

func (s *MultiwarehouseService) ListWarehouses(ctx context.Context, limit int) ([]*ent.Warehouse, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	ws, err := s.client.Warehouse.Query().
		Order(ent.Asc(warehouse.FieldCode)).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("list warehouses: %w", err)
	}
	return ws, nil
}

func (s *MultiwarehouseService) AssignLocation(ctx context.Context, whCode, locCode string) error {
	whCode = strings.TrimSpace(whCode)
	locCode = strings.TrimSpace(locCode)
	if whCode == "" {
		return ErrInvalidWarehouseCode
	}
	if locCode == "" {
		return ErrInvalidLocationCode
	}

	w, err := s.client.Warehouse.Query().
		Where(warehouse.Code(whCode)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return ErrWarehouseNotFound
		}
		return fmt.Errorf("fetch warehouse: %w", err)
	}

	loc, err := s.client.Location.Query().
		Where(location.Code(locCode)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return ErrLocationNotFound
		}
		return fmt.Errorf("fetch location: %w", err)
	}

	assigned, err := s.client.WarehouseLocation.Query().
		Where(warehouselocation.HasLocationWith(location.ID(loc.ID))).
		Exist(ctx)
	if err != nil {
		return fmt.Errorf("check assginment: %w", err)
	}
	if assigned {
		return ErrLocationAlreadyAssigned
	}

	_, err = s.client.WarehouseLocation.Create().
		SetWarehouse(w).
		SetLocation(loc).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return ErrLocationAlreadyAssigned
		}
		return fmt.Errorf("assgin location: %w", err)
	}
	return nil
}

func (s *MultiwarehouseService) ListLocations(ctx context.Context, whCode string, limit int) ([]*ent.Location, error) {
	whCode = strings.TrimSpace(whCode)
	if whCode == "" {
		return nil, ErrInvalidWarehouseCode
	}
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	w, err := s.client.Warehouse.Query().
		Where(warehouse.Code(whCode)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrWarehouseNotFound
		}
		return nil, fmt.Errorf("fetch warehouse: %w", err)
	}

	links, err := s.client.WarehouseLocation.Query().
		Where(warehouselocation.HasWarehouseWith(warehouse.ID(w.ID))).
		WithLocation().
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("list warehouse locations: %w", err)
	}

	out := make([]*ent.Location, 0, len(links))
	for _, l := range links {
		if l.Edges.Location != nil {
			out = append(out, l.Edges.Location)
		}
	}
	return out, nil
}
