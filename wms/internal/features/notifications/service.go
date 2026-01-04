//go:build notifications

package notifications

import (
	"context"
	"fmt"
	"time"
)

type NotificationService struct {
	cfg Config
}

func NewNotificationService(cfg Config) *NotificationService {
	return &NotificationService{cfg: cfg}
}

func (s *NotificationService) Config() Config {
	return s.cfg
}

func (s *NotificationService) Send(ctx context.Context, subject, message string) error {
	_ = ctx

	if !s.cfg.Enabled {
		return nil
	}

	ts := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] NOTIFY subject=%q recipients=%v message=%q\n", ts, subject, s.cfg.Recipients, message)
	return nil
}
