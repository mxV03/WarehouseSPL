//go:build reporting && multiwarehouse

package cli

import (
	"context"
	"fmt"

	"github.com/mxV03/wms/internal/features/interfaces/cli/clictx"
	"github.com/mxV03/wms/internal/features/interfaces/cli/registry"
	"github.com/mxV03/wms/internal/features/reporting"
)

func init() {
	registry.Register(registry.Command{
		Name:        "reporting.warehouse.summary",
		Usage:       "reporting.warehouse.summary <warehouseCode>",
		Group:       "Optional / Reporting",
		Description: "Warehouse summary KPIs (locations, stock movement totals).",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("usage: reporting.warehouse.summary <warehouseCode>")
			}

			client := clictx.AppCtx().Client()
			s, err := reporting.WarehouseSummaryReport(ctx, client, args[0])
			if err != nil {
				return err
			}

			fmt.Printf("Warehouse %s (%s)\n", s.Code, dash(s.Name))
			fmt.Printf("Locations: %d\n", s.LocationCount)
			fmt.Printf("Movements: %d\n", s.MovementCount)
			fmt.Printf("Total IN:  %d\n", s.TotalIn)
			fmt.Printf("Total OUT: %d\n", s.TotalOut)
			fmt.Printf("NET:       %d\n", s.Net)

			return nil
		},
	})
}

func dash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
