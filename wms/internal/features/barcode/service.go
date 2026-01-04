//go:build barcode

package barcode

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mxV03/wms/ent"
	"github.com/mxV03/wms/ent/bin"
	"github.com/mxV03/wms/ent/item"
	"github.com/mxV03/wms/ent/location"
)

type Config struct {
	Format     string
	ItemPrefix string
	BinPrefix  string
}

func LoadConfig() Config {
	c := Config{
		Format:     strings.TrimSpace(os.Getenv("WMS_BARCODE_FORMAT")),
		ItemPrefix: strings.TrimSpace(os.Getenv("WMS_BARCODE_ITEM_PREFIX")),
		BinPrefix:  strings.TrimSpace(os.Getenv("WMS_BARCODE_BIN_PREFIX")),
	}
	if c.Format == "" {
		c.Format = "plain"
	}
	if c.ItemPrefix == "" {
		c.ItemPrefix = "ITEM:"
	}
	if c.BinPrefix == "" {
		c.BinPrefix = "BIN:"
	}
	return c
}

type BarcodeService struct {
	client *ent.Client
	cfg    Config
}

type ScanResult struct {
	Kind       string
	SKU        string
	Location   string
	Bin        string
	ExistsInDB bool
}

func NewBarcodeService(client *ent.Client) *BarcodeService {
	return &BarcodeService{
		client: client,
		cfg:    LoadConfig(),
	}
}

func (s *BarcodeService) Config() Config {
	return s.cfg
}

func (s *BarcodeService) PrintItemBarcode(ctx context.Context, sku string) (string, error) {
	sku = strings.TrimSpace(sku)
	if sku == "" {
		return "", fmt.Errorf("invalid sku")
	}

	if _, err := s.client.Item.Query().Where(item.SKU(sku)).Only(ctx); err != nil {
		if ent.IsNotFound(err) {
			return "", fmt.Errorf("item not found")
		}
		return "", fmt.Errorf("fetch item: %w", err)
	}
	return s.encode(s.cfg.ItemPrefix + sku), nil
}

func (s *BarcodeService) PrintBinBarcode(ctx context.Context, locCode, binCode string) (string, error) {
	locCode = strings.TrimSpace(locCode)
	binCode = strings.TrimSpace(binCode)
	if locCode == "" {
		return "", fmt.Errorf("invalid location code")
	}
	if binCode == "" {
		return "", fmt.Errorf("invalid bin code")
	}

	loc, err := s.client.Location.Query().Where(location.Code(locCode)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return "", fmt.Errorf("location not found")
		}
		return "", fmt.Errorf("fetch location: %w", err)
	}

	if _, err := s.client.Bin.Query().Where(bin.Code(binCode), bin.HasLocationWith(location.ID(loc.ID))).Only(ctx); err != nil {
		if ent.IsNotFound(err) {
			return "", fmt.Errorf("bin not found for location")
		}
		return "", fmt.Errorf("fetch bin: %w", err)
	}

	payload := fmt.Sprintf("%s|%s", locCode, binCode)
	return s.encode(s.cfg.BinPrefix + payload), nil
}

func (s *BarcodeService) Scan(ctx context.Context, code string) (ScanResult, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return ScanResult{}, fmt.Errorf("empty code")
	}

	raw := s.decode(code)

	if strings.HasPrefix(raw, s.cfg.ItemPrefix) {
		sku := strings.TrimPrefix(raw, s.cfg.ItemPrefix)
		sku = strings.TrimSpace(sku)
		if sku == "" {
			return ScanResult{}, fmt.Errorf("invalid item barcode")
		}

		exists, err := s.client.Item.Query().
			Where(item.SKU(sku)).
			Exist(ctx)
		if err != nil {
			return ScanResult{}, fmt.Errorf("check item: %w", err)
		}

		return ScanResult{
			Kind:       "ITEM",
			SKU:        sku,
			ExistsInDB: exists,
		}, nil
	}

	if strings.HasPrefix(raw, s.cfg.BinPrefix) {
		payload := strings.TrimPrefix(raw, s.cfg.BinPrefix)
		payload = strings.TrimSpace(payload)
		parts := strings.Split(payload, "|")
		if len(parts) != 2 {
			return ScanResult{}, fmt.Errorf("invalid bin barcode")
		}
		locCode := strings.TrimSpace(parts[0])
		binCode := strings.TrimSpace(parts[1])
		if locCode == "" || binCode == "" {
			return ScanResult{}, fmt.Errorf("invalid bin barcode")
		}

		loc, err := s.client.Location.Query().
			Where(location.Code(locCode)).
			Only(ctx)

		if err != nil {
			if ent.IsNotFound(err) {
				return ScanResult{Kind: "BIN", Location: locCode, Bin: binCode, ExistsInDB: false}, nil
			}
			return ScanResult{}, fmt.Errorf("fetch location: %w", err)
		}

		exists, err := s.client.Bin.Query().
			Where(bin.Code(binCode), bin.HasLocationWith(location.ID(loc.ID))).
			Exist(ctx)
		if err != nil {
			return ScanResult{}, fmt.Errorf("check bin: %w", err)
		}
		return ScanResult{
			Kind:       "BIN",
			Location:   locCode,
			Bin:        binCode,
			ExistsInDB: exists,
		}, nil
	}

	return ScanResult{Kind: "UNKNOWN"}, nil
}

func (s *BarcodeService) encode(payload string) string {
	switch strings.ToLower(s.cfg.Format) {
	case "plain":
		return payload
	case "bracket":
		return "[" + payload + "]"
	default:
		return payload
	}
}

func (s *BarcodeService) decode(code string) string {
	code = strings.TrimSpace(code)
	if strings.HasPrefix(code, "[") && strings.HasSuffix(code, "]") {
		return strings.TrimSuffix(strings.TrimPrefix(code, "["), "]")
	}
	return code
}
