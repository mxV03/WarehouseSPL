//go:build barcode

package cli

import (
	"context"
	"fmt"

	"github.com/mxV03/wms/internal/features/barcode"
	"github.com/mxV03/wms/internal/features/interfaces/cli/clictx"
	"github.com/mxV03/wms/internal/features/interfaces/cli/registry"
)

func init() {
	registry.Register(registry.Command{
		Name:        "barcode.config",
		Usage:       "barcode.config",
		Group:       "Optional / Barcode",
		Description: "Show runtime barcode configuration.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("usage: barcode.config")
			}
			svc := barcode.NewBarcodeService(clictx.AppCtx().Client())
			c := svc.Config()
			fmt.Printf("format=%s\nitem_prefix=%s\nbin_prefix=%s\n", c.Format, c.ItemPrefix, c.BinPrefix)
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "barcode.print.item",
		Usage:       "barcode.print.item <sku>",
		Group:       "Optional / Barcode",
		Description: "Print a mock barcode for an item SKU.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("usage: barcode.print.item <sku>")
			}
			svc := barcode.NewBarcodeService(clictx.AppCtx().Client())
			code, err := svc.PrintItemBarcode(ctx, args[0])
			if err != nil {
				return err
			}
			fmt.Println(code)
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "barcode.print.bin",
		Usage:       "barcode.print.bin <locationCode> <binCode>",
		Group:       "Optional / Barcode",
		Description: "Print a mock barcode for a bin (scoped by location).",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("usage: barcode.print.item <locationCode> <binCode>")
			}
			svc := barcode.NewBarcodeService(clictx.AppCtx().Client())
			code, err := svc.PrintBinBarcode(ctx, args[0], args[1])
			if err != nil {
				return err
			}
			fmt.Println(code)
			return nil
		},
	})

	registry.Register(registry.Command{
		Name:        "barcode.scan",
		Usage:       "barcode.scan <code>",
		Group:       "Optional / Barcode",
		Description: "Scan (decode) a barcode and optionally validate against the DB.",
		Run: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("usage: barcode.scan <code>")
			}
			svc := barcode.NewBarcodeService(clictx.AppCtx().Client())
			res, err := svc.Scan(ctx, args[0])
			if err != nil {
				return err
			}

			switch res.Kind {
			case "ITEM":
				fmt.Printf("Item sku=%s exists=%v\n", res.SKU, res.ExistsInDB)
			case "BIN":
				fmt.Printf("BIN loc=%s bin=%s exists=%v\n", res.Location, res.Bin, res.ExistsInDB)
			default:
				fmt.Println("UNKNOWN")
			}
			return nil
		},
	})
}
