//go:build reporting && multiwarehouse

package reporting

import (
	"context"
	"fmt"

	"github.com/mxV03/wms/ent"
	"github.com/mxV03/wms/ent/location"
	"github.com/mxV03/wms/ent/stockmovement"
	"github.com/mxV03/wms/ent/warehouse"
	"github.com/mxV03/wms/ent/warehouselocation"
)

type WarehouseSummary struct {
	Code          string
	Name          string
	LocationCount int
	MovementCount int
	TotalIn       int
	TotalOut      int
	Net           int
}

func WarehouseSummaryReport(ctx context.Context, client *ent.Client, whCode string) (*WarehouseSummary, error) {
	if whCode == "" {
		return nil, fmt.Errorf("invalid warehouse code")
	}

	w, err := client.Warehouse.Query().
		Where(warehouse.Code(whCode)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("warehouse not found")
		}
		return nil, fmt.Errorf("fetch warehouse: %w", err)
	}

	links, err := client.WarehouseLocation.Query().
		Where(warehouselocation.HasWarehouseWith(warehouse.ID(w.ID))).
		WithLocation().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch warehouse locations: %w", err)
	}

	locIDs := make([]int, 0, len(links))
	for _, l := range links {
		if l.Edges.Location != nil {
			locIDs = append(locIDs, l.Edges.Location.ID)
		}
	}

	summary := &WarehouseSummary{
		Code:          w.Code,
		Name:          w.Name,
		LocationCount: len(locIDs),
	}

	if len(locIDs) == 0 {
		return summary, nil
	}

	movs, err := client.StockMovement.Query().
		Where(stockmovement.HasLocationWith(location.IDIn(locIDs...))).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch stock movements: %w", err)
	}

	summary.MovementCount = len(movs)
	for _, m := range movs {
		switch m.Type {
		case "IN":
			summary.TotalIn += m.Quantity
		case "OUT":
			summary.TotalOut += m.Quantity
		}
	}
	summary.Net = summary.TotalIn - summary.TotalOut
	return summary, nil
}
