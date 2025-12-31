package cli

import (
	"context"
	"fmt"
	"strconv"

	corestock "github.com/mxV03/wms/internal/core/inventory/stock"
	"github.com/mxV03/wms/internal/features/interfaces/cli/clictx"
	"github.com/mxV03/wms/internal/features/interfaces/cli/registry"
)

func init() {
	registry.Register(registry.Command{
		Name:        "stock.in",
		Usage:       "stock.in <sku> <location_code> <quantity> [reference]",
		Group:       "Core / Stock",
		Description: "Book incoming stock.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) < 3 {
				return fmt.Errorf("usage: stock.in <sku> <location_code> <quantity> [reference]")
			}
			qty, err := strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("invalid quantity: %w", err)
			}
			ref := ""
			if len(args) >= 4 {
				ref = args[3]
			}

			stockService := corestock.NewStockService(clictx.AppCtx().Client())
			return stockService.IN(ctx, args[0], args[1], qty, ref)
		},
	})

	registry.Register(registry.Command{
		Name:        "stock.out",
		Usage:       "stock.out <sku> <location_code> <quantity> [reference]",
		Group:       "Core / Stock",
		Description: "Book outgoing stock.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) < 3 {
				return fmt.Errorf("usage: stock.out <sku> <location_code> <quantity> [reference]")
			}

			sku := args[0]
			locCode := args[1]

			qty, err := strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("invalid quantity: %w", err)
			}

			ref := ""
			if len(args) >= 4 {
				ref = args[3]
			}

			stockService := corestock.NewStockService(clictx.AppCtx().Client())
			return stockService.OUT(ctx, sku, locCode, qty, ref)
		},
	})

	registry.Register(registry.Command{
		Name:        "stock.at",
		Usage:       "stock.at <sku> <location_code>",
		Group:       "Core / Stock",
		Description: "Show stock for a SKU at a specific location",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("usage: stock.at <sku> <location_code>")
			}

			svc := corestock.NewStockService(clictx.AppCtx().Client())
			qty, err := svc.StockAtLocation(ctx, args[0], args[1])
			if err != nil {
				return err
			}

			fmt.Printf("Stock for SKU=%s at location=%s; Quantity: %d\n", args[0], args[1], qty)
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "stock.total",
		Usage:       "stock.total <sku>",
		Group:       "Core / Stock",
		Description: "Show total stock for a SKU across all locations",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("usage: stock.total <sku>")
			}

			svc := corestock.NewStockService(clictx.AppCtx().Client())
			total, err := svc.StockBySKU(ctx, args[0])
			if err != nil {
				return err
			}

			fmt.Printf("Total stock for SKU=%s; Quantity: %d\n", args[0], total)
			return nil
		},
	})
}
