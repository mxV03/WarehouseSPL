package stock

import (
	"context"
	"fmt"
	"strings"

	"github.com/mxV03/wms/ent"
	"github.com/mxV03/wms/ent/item"
	"github.com/mxV03/wms/ent/location"
	"github.com/mxV03/wms/ent/stockmovement"
)

var (
	ErrInvalidQuantity   = fmt.Errorf("invalid stock quantity")
	ErrInsufficientStock = fmt.Errorf("insufficient stock available")
	ErrInvalidSKU        = fmt.Errorf("invalid stock SKU")
	ErrInvalidLocation   = fmt.Errorf("invalid stock location")
)

type MovementType string

const (
	MovementTypeIn  MovementType = "IN"
	MovementTypeOut MovementType = "OUT"
	// MoveMove MovementType = "MOVE"
)

type StockService struct {
	client *ent.Client
}

func NewStockService(client *ent.Client) *StockService {
	return &StockService{client: client}
}

type StockDTO struct {
	SKU          string
	LocationCode string
	Quantity     int
}

func (s *StockService) IN(ctx context.Context, sku, locCode string, qty int, ref string) error {
	sku = strings.TrimSpace(sku)
	locCode = strings.TrimSpace(locCode)
	ref = strings.TrimSpace(ref)

	if sku == "" {
		return ErrInvalidSKU
	}
	if locCode == "" {
		return ErrInvalidLocation
	}
	if qty <= 0 {
		return ErrInvalidQuantity
	}

	item, err := s.client.Item.Query().Where(item.SKU(sku)).Only(ctx)
	if err != nil {
		return fmt.Errorf("fetching item: %w", err)
	}

	loc, err := s.client.Location.Query().Where(location.Code(locCode)).Only(ctx)
	if err != nil {
		return fmt.Errorf("fetching location: %w", err)
	}

	create := s.client.StockMovement.Create().
		SetItem(item).
		SetLocation(loc).
		SetQuantity(qty).
		SetType(string(MovementTypeIn))

	if ref != "" {
		create.SetReference(ref)
	}

	if _, err := create.Save(ctx); err != nil {
		return fmt.Errorf("creating stock movement IN: %w", err)
	}
	return nil
}

func (s *StockService) OUT(ctx context.Context, sku, locCode string, qty int, ref string) error {
	sku = strings.TrimSpace(sku)
	locCode = strings.TrimSpace(locCode)
	ref = strings.TrimSpace(ref)

	if sku == "" {
		return ErrInvalidSKU
	}
	if locCode == "" {
		return ErrInvalidLocation
	}
	if qty <= 0 {
		return ErrInvalidQuantity
	}

	current, err := s.StockAtLocation(ctx, sku, locCode)
	if err != nil {
		return err
	}
	if current < qty {
		return ErrInsufficientStock
	}

	item, err := s.client.Item.Query().Where(item.SKU(sku)).Only(ctx)
	if err != nil {
		return fmt.Errorf("fetching item: %w", err)
	}

	loc, err := s.client.Location.Query().Where(location.Code(locCode)).Only(ctx)
	if err != nil {
		return fmt.Errorf("fetching location: %w", err)
	}

	create := s.client.StockMovement.Create().
		SetItem(item).
		SetLocation(loc).
		SetQuantity(qty).
		SetType(string(MovementTypeOut))

	if ref != "" {
		create.SetReference(ref)
	}

	if _, err := create.Save(ctx); err != nil {
		return fmt.Errorf("creating stock movement OUT: %w", err)
	}
	return nil
}

func (s *StockService) StockAtLocation(ctx context.Context, sku, locCode string) (int, error) {
	sku = strings.TrimSpace(sku)
	locCode = strings.TrimSpace(locCode)

	if sku == "" {
		return 0, ErrInvalidSKU
	}
	if locCode == "" {
		return 0, ErrInvalidLocation
	}

	q := s.client.StockMovement.Query().
		Where(
			stockmovement.HasItemWith(item.SKU(sku)),
			stockmovement.HasLocationWith(location.Code(locCode)),
		)

	ins, err := q.Clone().
		Where(stockmovement.TypeEQ(string(MovementTypeIn))).
		Aggregate(ent.Sum(stockmovement.FieldQuantity)).
		Int(ctx)

	if err != nil {
		return 0, fmt.Errorf("aggregating IN stock: %w", err)
	}

	outs, err := q.Clone().
		Where(stockmovement.TypeEQ(string(MovementTypeOut))).
		Aggregate(ent.Sum(stockmovement.FieldQuantity)).
		Int(ctx)

	if err != nil {
		return 0, fmt.Errorf("aggregating OUT stock: %w", err)
	}

	return ins - outs, nil
}

// total stock of an item across all locations
func (s *StockService) StockBySKU(ctx context.Context, sku string) (int, error) {
	sku = strings.TrimSpace(sku)
	if sku == "" {
		return 0, ErrInvalidSKU
	}

	q := s.client.StockMovement.Query().
		Where(
			stockmovement.HasItemWith(item.SKU(sku)),
		)
	ins, err := q.Clone().
		Where(stockmovement.TypeEQ(string(MovementTypeIn))).
		Aggregate(ent.Sum(stockmovement.FieldQuantity)).
		Int(ctx)
	if err != nil {
		return 0, fmt.Errorf("aggregating IN stock: %w", err)
	}

	outs, err := q.Clone().
		Where(stockmovement.TypeEQ(string(MovementTypeOut))).
		Aggregate(ent.Sum(stockmovement.FieldQuantity)).
		Int(ctx)
	if err != nil {
		return 0, fmt.Errorf("aggregating OUT stock: %w", err)
	}

	return ins - outs, nil
}
