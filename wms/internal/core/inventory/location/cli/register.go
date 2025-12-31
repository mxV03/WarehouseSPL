package cli

import (
	"context"
	"fmt"
	"strconv"

	corelocation "github.com/mxV03/wms/internal/core/inventory/location"
	"github.com/mxV03/wms/internal/features/interfaces/cli/clictx"
	"github.com/mxV03/wms/internal/features/interfaces/cli/registry"
)

func init() {
	registry.Register(registry.Command{
		Name:        "location.add",
		Usage:       "location.add <code> <name>",
		Group:       "Core / Location",
		Description: "Create a new location.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) < 2 {
				return fmt.Errorf("usage: location.add <code> <name>")
			}

			svc := corelocation.NewLocationService(clictx.AppCtx().Client())
			dto, err := svc.CreateLocation(ctx, args[0], args[1])
			if err != nil {
				return err
			}
			fmt.Printf("created location: CODE=%s ID=%d\n", dto.Code, dto.ID)
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "location.get",
		Usage:       "location.get <code>",
		Group:       "Core / Location",
		Description: "Get a location by code.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("usage: location.get <code>")
			}

			svc := corelocation.NewLocationService(clictx.AppCtx().Client())
			dto, err := svc.GetLocationByCode(ctx, args[0])
			if err != nil {
				return err
			}
			fmt.Printf("location: CODE=%s NAME=%s\n", dto.Code, dto.Name)
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "location.list",
		Usage:       "location.list [limit]",
		Group:       "Core / Location",
		Description: "List locations. (default limit=100, max=500)",
		Run: func(ctx context.Context, args []string) error {
			limit := 100
			if len(args) == 1 {
				v, err := strconv.Atoi(args[0])
				if err != nil {
					return fmt.Errorf("limit must be an integer")
				}
				limit = v
			} else if len(args) > 1 {
				return fmt.Errorf("usage: location.list [limit]")
			}

			svc := corelocation.NewLocationService(clictx.AppCtx().Client())
			locations, err := svc.ListLocations(ctx, limit)
			if err != nil {
				return err
			}

			if len(locations) == 0 {
				fmt.Println("no locations found")
				return nil
			}

			for _, dto := range locations {
				fmt.Printf("location: CODE=%s NAME=%s\n", dto.Code, dto.Name)
			}
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "location.del",
		Usage:       "location.del <code>",
		Group:       "Core / Location",
		Description: "Delete a location by code.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("usage: location.del <code>")
			}

			svc := corelocation.NewLocationService(clictx.AppCtx().Client())
			if err := svc.DeleteLocationByCode(ctx, args[0]); err != nil {
				return err
			}

			fmt.Printf("deleted location with CODE=%s\n", args[0])
			return nil
		},
	})
}
