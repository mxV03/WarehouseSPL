//go:build multiwarehouse

package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/mxV03/wms/internal/features/interfaces/cli/clictx"
	"github.com/mxV03/wms/internal/features/interfaces/cli/registry"
	"github.com/mxV03/wms/internal/features/multiwarehouse"
)

func init() {
	registry.Register(registry.Command{
		Name:        "warehouse.add",
		Usage:       "warehouse.add <warehouseCode> [name]",
		Group:       "Optional / MultiWarehouse",
		Description: "Create a warehouse.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) < 1 || len(args) > 2 {
				return fmt.Errorf("usage: warehouse.add <warehouseCode> [name]")
			}
			name := ""
			if len(args) == 2 {
				name = args[1]
			}
			svc := multiwarehouse.NewMultiwarehouseService(clictx.AppCtx().Client())
			w, err := svc.CreateWarehouse(ctx, args[0], name)
			if err != nil {
				return err
			}
			fmt.Printf("created warehouse: CODE=%s NAME=%s\n", w.Code, empty(w.Name))
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "warehouse.list",
		Usage:       "warehouse.lsit [limit]",
		Group:       "Optional / MultiWarehouse",
		Description: "List warehouses.",
		Run: func(ctx context.Context, args []string) error {
			limit := 100

			if len(args) > 1 {
				return fmt.Errorf("usage: warehouse.lsit [limit]")
			}
			if len(args) == 1 {
				v, err := strconv.Atoi(args[0])
				if err != nil {
					return fmt.Errorf("limit must be an integer")
				}
				limit = v
			}
			svc := multiwarehouse.NewMultiwarehouseService(clictx.AppCtx().Client())
			ws, err := svc.ListWarehouses(ctx, limit)
			if err != nil {
				return err
			}
			if len(ws) == 0 {
				fmt.Println("no warehouses")
				return nil
			}
			for _, w := range ws {
				fmt.Printf("warehouse: CODE=%s NAME=%s\n", w.Code, empty(w.Name))
			}
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "warehouse.location.assign",
		Usage:       "warehouse.location.assign <warehouseCode> <locationCode>",
		Group:       "Optional / MultiWarehouse",
		Description: "Assign an existing location to a warehouse.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("usage: warehouse.location.assign <warehouseCode> <locationCode>")
			}
			svc := multiwarehouse.NewMultiwarehouseService(clictx.AppCtx().Client())
			if err := svc.AssignLocation(ctx, args[0], args[1]); err != nil {
				return err
			}
			fmt.Printf("assigned location %s to warehouse %s\n", args[1], args[0])
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "warehouse.location.list",
		Usage:       "warehouse.location.list <warehouseCode> [limit]",
		Group:       "Optional / MultiWarehouse",
		Description: "List locations assigned to a warehouse.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) < 1 || len(args) > 2 {
				return fmt.Errorf("usage: warehouse.location.list <warehouseCode> [limit]")
			}
			limit := 100
			if len(args) == 2 {
				v, err := strconv.Atoi(args[1])
				if err != nil {
					return fmt.Errorf("limit must be an integer")
				}
				limit = v
			}

			svc := multiwarehouse.NewMultiwarehouseService(clictx.AppCtx().Client())
			locs, err := svc.ListLocations(ctx, args[0], limit)
			if err != nil {
				return err
			}
			if len(locs) == 0 {
				fmt.Println("no locations assigned")
				return nil
			}
			for _, l := range locs {
				fmt.Printf("location: CODE=%s NAME=%s\n", l.Code, empty(l.Name))
			}
			return nil
		},
	})
}

func empty(s string) string {
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return s
}
