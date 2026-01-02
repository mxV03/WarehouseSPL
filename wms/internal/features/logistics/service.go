//go:build logistics

package logistics

import (
	"context"
	"fmt"
	"strings"

	"github.com/mxV03/wms/ent"
	"github.com/mxV03/wms/ent/bin"
	"github.com/mxV03/wms/ent/item"
	"github.com/mxV03/wms/ent/location"
	"github.com/mxV03/wms/ent/zone"
)

var (
	ErrInvalidLocation  = fmt.Errorf("invalid location code")
	ErrInvalidZonesCode = fmt.Errorf("invalid zone code")
	ErrInvalidBinCode   = fmt.Errorf("invalid bin code")
	ErrInvalidSKU       = fmt.Errorf("invalid sku")
	ErrNotFound         = fmt.Errorf("not found")
)

type LogisticsService struct {
	client *ent.Client
}

func NewLocationService(client *ent.Client) *LogisticsService {
	return &LogisticsService{
		client: client,
	}
}

func (s *LogisticsService) getLocation(ctx context.Context, locCode string) (*ent.Location, error) {
	locCode = strings.TrimSpace(locCode)
	if locCode == "" {
		return nil, ErrInvalidLocation
	}

	loc, err := s.client.Location.Query().
		Where(location.Code(locCode)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("fetch location: %w", err)
	}
	return loc, nil
}

func (s *LogisticsService) CreateZone(ctx context.Context, locCode, zoneCode, name string) error {
	zoneCode = strings.TrimSpace(zoneCode)
	name = strings.TrimSpace(name)
	if zoneCode == "" {
		return ErrInvalidZonesCode
	}

	loc, err := s.getLocation(ctx, locCode)
	if err != nil {
		return err
	}

	_, err = s.client.Zone.Create().
		SetCode(zoneCode).
		SetName(name).
		SetLocation(loc).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("create zone: %w", err)
	}

	return nil
}

func (s *LogisticsService) ListZones(ctx context.Context, locCode string, limit int) ([]*ent.Zone, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	loc, err := s.getLocation(ctx, locCode)
	if err != nil {
		return nil, err
	}
	zs, err := s.client.Zone.Query().
		Where(zone.HasLocationWith(location.ID(loc.ID))).
		Limit(limit).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("list zones: %w", err)
	}
	return zs, nil
}

func (s *LogisticsService) CreateBin(ctx context.Context, locCode, zoneCode, binCode, name string) error {
	binCode = strings.TrimSpace(binCode)
	zoneCode = strings.TrimSpace(zoneCode)
	name = strings.TrimSpace(name)
	if binCode == "" {
		return ErrInvalidBinCode
	}
	if zoneCode == "" {
		return ErrInvalidZonesCode
	}

	loc, err := s.getLocation(ctx, locCode)
	if err != nil {
		return err
	}

	z, err := s.client.Zone.Query().
		Where(zone.Code(zoneCode), zone.HasLocationWith(location.ID(loc.ID))).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("zone not found for location")
		}
		return fmt.Errorf("fetch zone: %w", err)
	}

	_, err = s.client.Bin.Create().
		SetCode(binCode).
		SetName(name).
		SetLocation(loc).
		SetZone(z).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("create bin: %w", err)
	}
	return nil
}

func (s *LogisticsService) ListBins(ctx context.Context, locCode string, zoneCode string, limit int) ([]*ent.Bin, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	loc, err := s.getLocation(ctx, locCode)
	if err != nil {
		return nil, err
	}

	q := s.client.Bin.Query().
		Where(bin.HasLocationWith(location.ID(loc.ID)))

	zoneCode = strings.TrimSpace(zoneCode)
	if zoneCode != "" {
		q = q.Where(bin.HasZoneWith(zone.Code(zoneCode)))
	}

	bs, err := q.Limit(limit).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("list bins: %w", err)
	}

	return bs, nil
}

func (s *LogisticsService) AssignItemToBin(ctx context.Context, locCode, binCode, sku string) error {
	binCode = strings.TrimSpace(binCode)
	sku = strings.TrimSpace(sku)
	if binCode == "" {
		return ErrInvalidBinCode
	}
	if sku == "" {
		return ErrInvalidSKU
	}

	loc, err := s.getLocation(ctx, locCode)
	if err != nil {
		return err
	}

	b, err := s.client.Bin.Query().
		Where(bin.Code(binCode), bin.HasLocationWith(location.ID(loc.ID))).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("fetch bin: %w", err)
	}

	it, err := s.client.Item.Query().
		Where(item.SKU(sku)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("fetch item: %w", err)
	}

	if err := b.Update().AddItems(it).Exec(ctx); err != nil {
		return fmt.Errorf("assgin item to bin: %w", err)
	}

	return nil
}

func (s *LogisticsService) ItemsInBin(ctx context.Context, locCode, binCode string, limit int) ([]*ent.Item, error) {
	binCode = strings.TrimSpace(binCode)
	if binCode == "" {
		return nil, ErrInvalidBinCode
	}
	if limit <= 0 || limit > 500 {
		limit = 100
	}

	loc, err := s.getLocation(ctx, locCode)
	if err != nil {
		return nil, err
	}

	b, err := s.client.Bin.Query().
		Where(bin.Code(binCode), bin.HasLocationWith(location.ID(loc.ID))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("fetch bin: %w", err)
	}

	items, err := b.QueryItems().
		Order(ent.Asc(item.FieldSKU)).
		Limit(limit).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("list items in bin: %w", err)
	}
	return items, nil
}

func (s *LogisticsService) UnassignItemFromBin(ctx context.Context, locCode, binCode, sku string) error {
	binCode = strings.TrimSpace(binCode)
	sku = strings.TrimSpace(sku)

	if binCode == "" {
		return ErrInvalidBinCode
	}
	if sku == "" {
		return ErrInvalidSKU
	}

	loc, err := s.getLocation(ctx, locCode)
	if err != nil {
		return err
	}

	b, err := s.client.Bin.Query().
		Where(bin.Code(binCode), bin.HasLocationWith(location.ID(loc.ID))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("fetch bin: %w", err)
	}

	it, err := s.client.Item.Query().
		Where(item.SKU(sku)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return ErrNotFound
		}
		fmt.Errorf("fetch item: %w", err)
	}

	if err := b.Update().RemoveItems(it).Exec(ctx); err != nil {
		return fmt.Errorf("unassign item from bin: %w", err)
	}
	return nil
}

func (s *LogisticsService) DeleteBin(ctx context.Context, locCode, binCode string) error {
	binCode = strings.TrimSpace(binCode)
	if binCode == "" {
		return ErrInvalidBinCode
	}

	loc, err := s.getLocation(ctx, locCode)
	if err != nil {
		return err
	}

	b, err := s.client.Bin.Query().
		Where(bin.Code(binCode), bin.HasLocationWith(location.ID(loc.ID))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("fetch bin: %w", err)
	}

	count, err := b.QueryItems().Count(ctx)
	if err != nil {
		return fmt.Errorf("count bin items: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("cannot delte bin %s: %d item(s) still assigned", binCode, count)
	}

	if err := s.client.Bin.DeleteOne(b).Exec(ctx); err != nil {
		return fmt.Errorf("delte bin: %w", err)
	}
	return nil
}

func (s *LogisticsService) DeleteZone(ctx context.Context, locCode, zoneCode string) error {
	zoneCode = strings.TrimSpace(zoneCode)
	if zoneCode == "" {
		return ErrInvalidZonesCode
	}

	loc, err := s.getLocation(ctx, locCode)
	if err != nil {
		return err
	}

	z, err := s.client.Zone.Query().
		Where(zone.Code(zoneCode), zone.HasLocationWith(location.ID(loc.ID))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return ErrNotFound
		}
		return fmt.Errorf("fetch zone: %w", err)
	}

	binCount, err := z.QueryBins().Count(ctx)
	if err != nil {
		return fmt.Errorf("count zone bins: %w", err)
	}
	if binCount > 0 {
		return fmt.Errorf("cannot delte zone %s: %d bin(s) still exists", zoneCode, binCount)
	}

	if err := s.client.Zone.DeleteOne(z).Exec(ctx); err != nil {
		return fmt.Errorf("delete zone: %w", err)
	}
	return nil
}
