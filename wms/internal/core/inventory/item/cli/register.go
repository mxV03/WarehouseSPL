package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	coreitem "github.com/mxV03/wms/internal/core/inventory/item"
	"github.com/mxV03/wms/internal/features/interfaces/cli/clictx"
	"github.com/mxV03/wms/internal/features/interfaces/cli/registry"
)

func init() {
	registry.Register(registry.Command{
		Name:        "item.add",
		Usage:       "item.add <sku> <name> [description]",
		Group:       "Core / Items",
		Description: "Create a new item.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) < 2 {
				return fmt.Errorf("usage: item.add <sku> <name> [description]")
			}
			desc := ""
			if len(args) > 2 {
				desc = strings.Join(args[2:], " ")
			}
			svc := coreitem.NewItemService(clictx.AppCtx().Client())
			dto, err := svc.CreateItem(ctx, args[0], args[1], desc)
			if err != nil {
				return err
			}
			fmt.Printf("created item: SKU=%s ID=%d\n", dto.SKU, dto.ID)
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "item.get",
		Usage:       "item.get <sku>",
		Group:       "Core / Items",
		Description: "Get item by SKU.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("usage: item.get <sku>")
			}
			svc := coreitem.NewItemService(clictx.AppCtx().Client())
			dto, err := svc.GetItemBySKU(ctx, args[0])
			if err != nil {
				return err
			}
			fmt.Printf("item: SKU=%s NAME=%s DESC=%s\n", dto.SKU, dto.Name, dto.Description)
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "item.list",
		Usage:       "item.list [limit]",
		Group:       "Core / Items",
		Description: "List all items. (default limit=100, max=500)",
		Run: func(ctx context.Context, args []string) error {
			limit := 100
			if len(args) == 1 {
				v, err := strconv.Atoi(args[0])
				if err != nil {
					return fmt.Errorf("limit must be an integer")
				}
				limit = v
			} else if len(args) > 1 {
				return fmt.Errorf("usage: item.list [limit]")
			}

			svc := coreitem.NewItemService(clictx.AppCtx().Client())
			items, err := svc.ListItems(ctx, limit)
			if err != nil {
				return err
			}

			if len(items) == 0 {
				fmt.Println("no items found")
				return nil
			}

			for _, dto := range items {
				fmt.Printf("item: SKU=%s NAME=%s DESC=%s\n", dto.SKU, dto.Name, dto.Description)
			}
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "item.del",
		Usage:       "item.del <sku>",
		Group:       "Core / Items",
		Description: "Delete item by SKU.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("usage: item.del <sku>")
			}

			svc := coreitem.NewItemService(clictx.AppCtx().Client())
			if err := svc.DeleteItemBySKU(ctx, args[0]); err != nil {
				return err
			}

			fmt.Printf("deleted item with SKU=%s\n", args[0])
			return nil
		},
	})
}
