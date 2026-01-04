//go:build picking

package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mxV03/wms/internal/features/interfaces/cli/clictx"
	"github.com/mxV03/wms/internal/features/interfaces/cli/registry"
	"github.com/mxV03/wms/internal/features/picking"
)

func init() {
	registry.Register(registry.Command{
		Name:        "picking.picklist.create",
		Usage:       "picking.picklist.create <orderNr>",
		Group:       "Optional / Picking",
		Description: "Create a picklist for an order (one task per order line).",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("usage: picking.picklist.create <orderNr>")
			}
			svc := picking.NewPickingService(clictx.AppCtx().Client())
			pl, err := svc.CreatePickList(ctx, args[0])
			if err != nil {
				return err
			}
			fmt.Printf("created picklist: ID=%d ORDER=%s\n", pl.ID, args[0])
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "picking.picklist.start",
		Usage:       "picking.picklist.start <pickListID>",
		Group:       "Optional / Picking",
		Description: "Start a picklist.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("usage: picking.picklist.start <pickListID>")
			}
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("pickListID must be an integer")
			}
			svc := picking.NewPickingService(clictx.AppCtx().Client())
			if err := svc.StartPickList(ctx, id); err != nil {
				return err
			}
			fmt.Printf("picklist started")
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "picking.task.pick",
		Usage:       "picking.task.pick <taskID>",
		Group:       "Optional / Picking",
		Description: "Mark one pick task as picked.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("usage: picking.task.pick <taskID>")
			}
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("taskID must be an integer")
			}
			svc := picking.NewPickingService(clictx.AppCtx().Client())
			if err := svc.MarkTaskPicked(ctx, id); err != nil {
				return err
			}
			fmt.Printf("task picked")
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "picking.picklist.done",
		Usage:       "picking.picklist.done <pickListID>",
		Group:       "Optional / Picking",
		Description: "Finish a picklist.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("usage: picking.picklist.done <pickListID>")
			}
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("pickListID must be an integer")
			}
			svc := picking.NewPickingService(clictx.AppCtx().Client())
			if err := svc.DonePickList(ctx, id); err != nil {
				return err
			}
			fmt.Printf("picklist done")
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "picking.picklist.show",
		Usage:       "picking.picklist.show <pickListID>",
		Group:       "Optional / Picking",
		Description: "Show picklist details.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("usage: picking.picklist.show <pickListID>")
			}
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("pickListID must be an integer")
			}
			svc := picking.NewPickingService(clictx.AppCtx().Client())
			pl, err := svc.ShowPickList(ctx, id)
			if err != nil {
				return err
			}
			fmt.Printf("PickList ID=%d ORDER=%s STATUS=%s\n", pl.ID, pl.OrderNr, pl.Status)
			for _, t := range pl.Tasks {
				fmt.Printf("  Task %d: SKU=%s QTY=%d LOC=%s BIN=%s STATUS=%s\n",
					t.ID, t.SKU, t.Quantity, t.Location, t.Bin, t.Status)
			}
			return nil
		},
	})
}
