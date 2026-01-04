//go:build notifications

package notifications

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Enabled           bool
	Recipients        []string
	LowStockThreshold int
}

func LoadConfigFromEnv() Config {
	cfg := Config{
		Enabled:           true,
		Recipients:        []string{"warehouse@local"},
		LowStockThreshold: 10,
	}

	if v := strings.TrimSpace(os.Getenv("WMS_NOTIFICATIONS_ENABLED")); v != "" {
		b, err := strconv.ParseBool(v)
		if err == nil {
			cfg.Enabled = b
		}
	}

	if v := strings.TrimSpace(os.Getenv("WMS_NOTIFICATIONS_RECIPIENTS")); v != "" {
		parts := strings.Split(v, ",")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				out = append(out, p)
			}
		}
		if len(out) > 0 {
			cfg.Recipients = out
		}
	}

	if v := strings.TrimSpace(os.Getenv("WMS_LOW_STOCK_THRESHOLD")); v != "" {
		n, err := strconv.Atoi(v)
		if err == nil && n >= 0 {
			cfg.LowStockThreshold = n
		}
	}

	return cfg
}
