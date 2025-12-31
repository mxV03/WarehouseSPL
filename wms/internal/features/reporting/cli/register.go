//go:build reporting

package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mxV03/warhousemanagementsystem/internal/features/interfaces/cli/clictx"
	"github.com/mxV03/warhousemanagementsystem/internal/features/interfaces/cli/registry"
	"github.com/mxV03/warhousemanagementsystem/internal/features/reporting"
)

func init() {
	registry.Register(registry.Command{
		Name:        "report.stock",
		Group:       "Optional / Reporting",
		Usage:       "report.stock <sku>",
		Description: "Show total stock for a SKU (Reporting feature).",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("usage: report.stock <sku>")
			}
			svc := reporting.NewReportService(clictx.AppCtx().Client())
			total, err := svc.StockTotal(ctx, args[0])
			if err != nil {
				return err
			}
			fmt.Printf("report stock %s = %d\n", args[0], total)
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "report.movements",
		Group:       "Optional / Reporting",
		Usage:       "reports.movements <sku> [limit]",
		Description: "Show recent stock movements for a SKU (Reporting feature).",
		Run: func(ctx context.Context, args []string) error {
			if len(args) < 1 || len(args) > 2 {
				return fmt.Errorf("usage: report.movements <sku> [limit]")
			}
			limit := 20
			if len(args) == 2 {
				v, err := strconv.Atoi(args[1])
				if err != nil {
					return fmt.Errorf("limit must be an integer")
				}
				limit = v
			}

			svc := reporting.NewReportService(clictx.AppCtx().Client())
			moves, err := svc.RecentMovementsBySKU(ctx, args[0], limit)
			if err != nil {
				return err
			}
			if len(moves) == 0 {
				fmt.Println("no movements found")
				return nil
			}
			for _, m := range moves {
				ref := m.Reference
				if ref == "" {
					ref = "-"
				}
				fmt.Printf("%s  %-3s  %-10s  qty=%d  ref=%s\n",
					m.CreatedAt.Format("2006-01-02 15:04"),
					m.Type,
					m.LocationCode,
					m.Quantity,
					ref,
				)
			}
			return nil
		},
	})
}
