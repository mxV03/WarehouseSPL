package cli

import (
	"context"
	"fmt"
	"strconv"

	coreorders "github.com/mxV03/warhousemanagementsystem/internal/core/ordermanagement/orders"
	"github.com/mxV03/warhousemanagementsystem/internal/features/interfaces/cli/clictx"
	"github.com/mxV03/warhousemanagementsystem/internal/features/interfaces/cli/registry"
)

func init() {
	registry.Register(registry.Command{
		Name:        "order.in",
		Usage:       "order.in <order_number>",
		Description: "Create inbound order.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("usage: order.in <order_number>")
			}
			orderService := coreorders.NewOrderService(clictx.AppCtx().Client())
			_, err := orderService.CreateInboundOrder(ctx, args[0])
			if err == nil {
				fmt.Printf("Inbound order '%s' created successfully.\n", args[0])
			}
			return err
		},
	})

	registry.Register(registry.Command{
		Name:        "order.out",
		Usage:       "order.out <order_number>",
		Description: "Create outbound order.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("usage: order.out <order_number>")
			}
			orderService := coreorders.NewOrderService(clictx.AppCtx().Client())
			_, err := orderService.CreateOutboundOrder(ctx, args[0])
			if err == nil {
				fmt.Printf("Outbound order '%s' created successfully.\n", args[0])
			}
			return err
		},
	})

	registry.Register(registry.Command{
		Name:        "order.addline",
		Usage:       "order.addline <order_number> <sku> <location> <quantity>",
		Description: "Add a line item to an order.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 4 {
				return fmt.Errorf("usage: order.addline <order_number> <sku> <location> <quantity>")
			}
			qty, err := strconv.Atoi(args[3])
			if err != nil {
				return fmt.Errorf("invalid quantity: %w", err)
			}

			orderService := coreorders.NewOrderService(clictx.AppCtx().Client())
			_, err = orderService.AddLine(ctx, args[0], args[1], args[2], qty)
			return err
		},
	})

	registry.Register(registry.Command{
		Name:        "order.post",
		Usage:       "order.post <order_number>",
		Description: "Post an order: creates stock movements and marks order as POSTED.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("usage: order.post <order_number>")
			}
			orderService := coreorders.NewOrderService(clictx.AppCtx().Client())
			if err := orderService.PostOrder(ctx, args[0]); err != nil {
				return err
			}
			fmt.Printf("Order '%s' posted successfully.\n", args[0])
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "order.cancel",
		Usage:       "order.cancel <order_number>",
		Description: "Cancel a DRAFT order.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("usage: order.cancel <order_number>")
			}
			orderService := coreorders.NewOrderService(clictx.AppCtx().Client())
			if err := orderService.CancelOrder(ctx, args[0]); err != nil {
				return err
			}
			fmt.Printf("Order '%s' cancelled successfully.\n", args[0])
			return nil
		},
	})
}
