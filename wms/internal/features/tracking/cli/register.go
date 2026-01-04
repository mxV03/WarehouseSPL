//go:build tracking

package cli

import (
	"context"
	"fmt"

	"github.com/mxV03/wms/internal/features/interfaces/cli/clictx"
	"github.com/mxV03/wms/internal/features/interfaces/cli/registry"
	"github.com/mxV03/wms/internal/features/tracking"
)

func init() {
	registry.Register(registry.Command{
		Name:        "tracking.set",
		Usage:       "tracking.set <orderNr> <trackingId> [trackingUrl] [carrier]",
		Group:       "Optional / Tracking",
		Description: "Attach or update tracking information for an oder",
		Run: func(ctx context.Context, args []string) error {
			if len(args) < 2 || len(args) > 4 {
				return fmt.Errorf("usage: tracking.set <orderNr> <trackingId> [trackingUrl] [carrier]")
			}

			orderNr := args[0]
			trackingID := args[1]

			trackingURL := ""
			if len(args) >= 3 {
				trackingURL = args[2]
			}

			carrier := ""
			if len(args) == 4 {
				carrier = args[3]
			}

			svc := tracking.NewTrackingService(clictx.AppCtx().Client())
			if err := svc.Set(ctx, orderNr, trackingID, trackingURL, carrier); err != nil {
				return err
			}

			fmt.Printf("tracking set: ORDER=%s TRACKING=%s\n", orderNr, trackingID)
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "tracking.get",
		Usage:       "tracking.get <orderNo>",
		Group:       "Optional / Tracking",
		Description: "Show tracking information for an order.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("usage: tracking.get <orderNo>")
			}

			svc := tracking.NewTrackingService(clictx.AppCtx().Client())
			dto, err := svc.Get(ctx, args[0])
			if err != nil {
				return err
			}

			if dto.TrackingID == "" {
				fmt.Printf("order=%s tracking=<none>\n", dto.OrderNr)
				return nil
			}

			fmt.Printf(
				"order=%s tracking_id=%s url=%s carrier=%s\n",
				dto.OrderNr,
				dto.TrackingID,
				empty(dto.TrackingURL),
				empty(dto.Carrier),
			)
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "tracking.clear",
		Usage:       "tracking.clear <orderNo>",
		Group:       "Optional / Tracking",
		Description: "Remove tracking information from an order.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("usage: tracking.clear <orderNo>")
			}

			svc := tracking.NewTrackingService(clictx.AppCtx().Client())
			if err := svc.Clear(ctx, args[0]); err != nil {
				return err
			}

			fmt.Printf("tracking cleared: ORDER=%s\n", args[0])
			return nil
		},
	})
}

func empty(s string) string {
	if s == "" {
		return "-"
	}

	return s
}
