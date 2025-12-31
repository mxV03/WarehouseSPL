//go:build reporting

package reporting

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mxV03/warhousemanagementsystem/ent"
	"github.com/mxV03/warhousemanagementsystem/ent/item"
	"github.com/mxV03/warhousemanagementsystem/ent/stockmovement"
	"github.com/mxV03/warhousemanagementsystem/internal/core/inventory/stock"
)

type ReportService struct {
	client   *ent.Client
	stockSvc *stock.StockService
}

func NewReportService(client *ent.Client) *ReportService {
	return &ReportService{
		client:   client,
		stockSvc: stock.NewStockService(client),
	}
}

type MovementDTO struct {
	Type         string
	SKU          string
	LocationCode string
	Quantity     int
	Reference    string
	CreatedAt    time.Time
}

// total stock of an item across all locations
func (s *ReportService) StockTotal(ctx context.Context, sku string) (int, error) {
	return s.stockSvc.StockBySKU(ctx, sku)
}

func (s *ReportService) RecentMovementsBySKU(ctx context.Context, sku string, limit int) ([]MovementDTO, error) {
	sku = strings.TrimSpace(sku)
	if sku == "" {
		return nil, fmt.Errorf("invalid SKU")
	}
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	moves, err := s.client.StockMovement.Query().
		Where(stockmovement.HasItemWith(item.SKU(sku))).
		WithItem().
		WithLocation().
		Order(ent.Desc(stockmovement.FieldCreatedAt)).
		Limit(limit).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("query movements: %w", err)
	}

	out := make([]MovementDTO, 0, len(moves))
	for _, m := range moves {
		out = append(out, MovementDTO{
			Type:         m.Type,
			SKU:          m.Edges.Item.SKU,
			LocationCode: m.Edges.Location.Code,
			Quantity:     m.Quantity,
			Reference:    m.Reference,
			CreatedAt:    m.CreatedAt,
		})
	}
	return out, nil
}
