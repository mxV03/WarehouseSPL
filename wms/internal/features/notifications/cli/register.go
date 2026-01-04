//go:build notifications

package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/mxV03/wms/internal/features/interfaces/cli/clictx"
	"github.com/mxV03/wms/internal/features/interfaces/cli/registry"
	"github.com/mxV03/wms/internal/features/notifications"
	"github.com/mxV03/wms/internal/features/reporting"
)

func init() {
	registry.Register(registry.Command{
		Name:        "notify.test",
		Group:       "Optional / Notifications",
		Usage:       "notify.test <message>",
		Description: "Send a test notification (prints to stdout).",
		Run: func(ctx context.Context, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("usage: notify.test <message>")
			}

			cfg := notifications.LoadConfigFromEnv()
			svc := notifications.NewNotificationService(cfg)

			msg := strings.Join(args, " ")
			return svc.Send(ctx, "TEST", msg)
		},
	})

	registry.Register(registry.Command{
		Name:        "notify.config",
		Group:       "Optional / Notifications",
		Usage:       "notify.config",
		Description: "Print current notification config (from env).",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("usage: notify.config")
			}

			_ = ctx

			cfg := notifications.LoadConfigFromEnv()
			fmt.Printf("enabled=%v recipients=%v low_stock_threshold=%d\n",
				cfg.Enabled, cfg.Recipients, cfg.LowStockThreshold,
			)

			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "notify.lowstock",
		Group:       "Optional / Notifications",
		Usage:       "notify.lowstock <sku> [threshold]",
		Description: "Send a notification if total stock for SKU is below threshold.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) < 1 || len(args) > 2 {
				return fmt.Errorf("usage: notify.lowstock <sku> [threshold]")
			}

			sku := args[0]
			cfg := notifications.LoadConfigFromEnv()
			threshold := cfg.LowStockThreshold

			if len(args) == 2 {
				n, err := strconv.Atoi(args[1])
				if err != nil || n < 0 {
					return fmt.Errorf("threshold must be a non-negative integer")
				}
				threshold = n
			}

			rep := reporting.NewReportService(clictx.AppCtx().Client())
			total, err := rep.StockTotal(ctx, sku)
			if err != nil {
				return err
			}

			if total < threshold {
				svc := notifications.NewNotificationService(cfg)
				subject := "LOW_STOCK"
				msg := fmt.Sprintf("SKU %s stock=%d is below threshold=%d", sku, total, threshold)
				return svc.Send(ctx, subject, msg)
			}

			fmt.Printf("ok: SKU %s stock=%d (threshold=%d)\n", sku, total, threshold)
			return nil

		},
	})
}
