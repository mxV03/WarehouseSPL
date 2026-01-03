//go:build logistics

package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/mxV03/wms/internal/features/interfaces/cli/clictx"
	"github.com/mxV03/wms/internal/features/interfaces/cli/registry"
	"github.com/mxV03/wms/internal/features/logistics"
)

func init() {
	registry.Register(registry.Command{
		Name:        "logistics.zone.add",
		Usage:       "logistics.zone.add <locationCode> <zoneCode> [name]",
		Group:       "Optional / Logistics",
		Description: "Create a zone for a location",
		Run: func(ctx context.Context, args []string) error {
			if len(args) < 2 || len(args) > 3 {
				return fmt.Errorf("usage: logistics.zone.add <locationCode> <zoneCode> [name]")
			}

			name := ""
			if len(args) == 3 {
				name = args[2]
			}

			svc := logistics.NewLocationService(clictx.AppCtx().Client())
			err := svc.CreateZone(ctx, args[0], args[1], name)
			if err != nil {
				return err
			}
			fmt.Printf("created zone: LOC=%s ZONE=%s", args[0], args[1])
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "logistics.zone.list",
		Usage:       "logistics.zone.list <locationCode> [zoneCode] [limit]",
		Group:       "Optional / Logistics",
		Description: "List zones of a location.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) < 1 || len(args) > 2 {
				return fmt.Errorf("usage: logistics.zone.list <locationCode> [zoneCode] [limit]")
			}

			limit := 100
			if len(args) == 2 {
				v, err := strconv.Atoi(args[1])
				if err != nil {
					return fmt.Errorf("limit must be an integer")
				}
				limit = v
			}

			svc := logistics.NewLocationService(clictx.AppCtx().Client())
			zs, err := svc.ListZones(ctx, args[0], limit)
			if err != nil {
				return err
			}
			if len(zs) == 0 {
				fmt.Println("no zones found")
				return nil
			}
			for _, z := range zs {
				name := strings.TrimSpace(z.Name)
				if name == "" {
					name = "-"
				}
				fmt.Printf("zone: CODE=%s NAME=%s\n", z.Code, name)
			}
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "logistics.zone.delete",
		Usage:       "logistics.zone.delete <locationCode> <zoneCode>",
		Group:       "Optional / Logistics",
		Description: "Delete a zone only if it has no bins.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("usage: logistics.bin.delete <locationCode> <zoneCode>")
			}

			svc := logistics.NewLocationService(clictx.AppCtx().Client())
			if err := svc.DeleteZone(ctx, args[0], args[1]); err != nil {
				return err
			}
			fmt.Printf("deleted zone %s (loc=%s)\n", args[1], args[0])
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "logistics.bin.add",
		Usage:       "logistics.bin.add <locationCode> <zoneCode> <binCod> [name]",
		Group:       "Optional / Logistics",
		Description: "Create a bin (storage location) inside a zone",
		Run: func(ctx context.Context, args []string) error {
			if len(args) < 3 || len(args) > 4 {
				return fmt.Errorf("usage: logistics.bin.add <locationCode> <zoneCode> <binCod> [name]")
			}

			name := ""
			if len(args) == 4 {
				name = args[3]
			}

			svc := logistics.NewLocationService(clictx.AppCtx().Client())
			err := svc.CreateBin(ctx, args[0], args[1], args[2], name)
			if err != nil {
				return err
			}
			fmt.Printf("created zone: LOC=%s ZONE=%s BIN=%s", args[0], args[1], args[2])
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "logistics.bin.list",
		Usage:       "logistics.bin.list <locationCode> [zoneCode] [limit]",
		Group:       "Optional / Logistics",
		Description: "List bins for a location, optionally filtered by zone",
		Run: func(ctx context.Context, args []string) error {
			if len(args) < 1 || len(args) > 3 {
				return fmt.Errorf("usage: logistics.bin.list <locationCode> [zoneCode] [limit]")
			}

			zoneCode := ""
			limit := 100
			if len(args) >= 2 {
				v, err := strconv.Atoi(args[1])
				if err == nil {
					limit = v
				} else {
					zoneCode = args[1]
				}
			}

			if len(args) == 3 {
				v, err := strconv.Atoi(args[2])
				if err != nil {
					return fmt.Errorf("limit must be an integer")
				}
				limit = v
			}

			svc := logistics.NewLocationService(clictx.AppCtx().Client())
			bs, err := svc.ListBins(ctx, args[0], zoneCode, limit)
			if err != nil {
				return err
			}
			if len(bs) == 0 {
				fmt.Println("no bins found")
				return nil
			}
			for _, z := range bs {
				name := strings.TrimSpace(z.Name)
				if name == "" {
					name = "-"
				}
				fmt.Printf("bin: CODE=%s NAME=%s\n", z.Code, name)
			}
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "logistics.bin.assign",
		Usage:       "logistics.bin.add <locationCode> <binCod> <sku>",
		Group:       "Optional / Logistics",
		Description: "Assign (link) an item SKU to a bin (where it is stored)",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("usage: logistics.bin.add <locationCode> <binCod> <sku>")
			}

			svc := logistics.NewLocationService(clictx.AppCtx().Client())
			if err := svc.AssignItemToBin(ctx, args[0], args[1], args[2]); err != nil {
				return err
			}
			fmt.Printf("assigned item %s to bin %s (loc=%s)\n", args[2], args[1], args[0])
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "logistics.bin.items",
		Usage:       "logistics.bin.items <locationCode> <binCode> [limit]",
		Group:       "Optional / Logistics",
		Description: "List items assigned to a bin.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) < 2 || len(args) > 3 {
				return fmt.Errorf("usage: logistics.bin.items <locationCode> <binCode> [limit]")
			}

			limit := 100
			if len(args) == 3 {
				v, err := strconv.Atoi(args[2])
				if err != nil {
					return fmt.Errorf("limit must be an integer")
				}
				limit = v
			}

			svc := logistics.NewLocationService(clictx.AppCtx().Client())
			its, err := svc.ItemsInBin(ctx, args[0], args[1], limit)
			if err != nil {
				return err
			}
			if len(its) == 0 {
				fmt.Println("no items assigned to this bin")
				return nil
			}
			for _, it := range its {
				fmt.Printf("item: SKU=%s NAME=%s\n", it.SKU, it.Name)
			}
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "logistics.bin.unassign",
		Usage:       "logistics.bin.unassign <locationCode> <binCode> <sku>",
		Group:       "Optional / Logistics",
		Description: "Unassign (unlink) an item SKU from a bin.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("usage: logistics.bin.unassign <locationCode> <binCode> <sku>")
			}

			svc := logistics.NewLocationService(clictx.AppCtx().Client())
			if err := svc.UnassignItemFromBin(ctx, args[0], args[1], args[2]); err != nil {
				return err
			}
			fmt.Printf("unassigned item %s from bin %s (loc=%s)\n", args[2], args[1], args[0])
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "logistics.bin.delete",
		Usage:       "logistics.bin.delete <locationCode> <binCode>",
		Group:       "Optional / Logistics",
		Description: "Delete a bin only if no items are assigned.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("usage: logistics.bin.delete <locationCode> <binCode>")
			}

			svc := logistics.NewLocationService(clictx.AppCtx().Client())
			if err := svc.DeleteBin(ctx, args[0], args[1]); err != nil {
				return err
			}
			fmt.Printf("deleted bin %s (loc=%s)\n", args[1], args[0])
			return nil
		},
	})

}
