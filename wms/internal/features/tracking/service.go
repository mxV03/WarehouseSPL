//go:build tracking

package tracking

import (
	"context"
	"fmt"
	"strings"

	"github.com/mxV03/wms/ent"
	"github.com/mxV03/wms/ent/order"
	"github.com/mxV03/wms/ent/tracking"
)

var (
	ErrInvalidOrderNr  = fmt.Errorf("invalid order number")
	ErrOrderNotFound   = fmt.Errorf("order not found")
	ErrInvalidTracking = fmt.Errorf("invalid tracking id")
)

type TrackingService struct {
	client *ent.Client
}

func NewTrackingService(client *ent.Client) *TrackingService {
	return &TrackingService{
		client: client,
	}
}

type TrackingDTO struct {
	OrderNr     string
	TrackingID  string
	TrackingURL string
	Carrier     string
}

func (s *TrackingService) Set(ctx context.Context, orderNr, trackingID, trackingURL, carrier string) error {
	orderNr = strings.TrimSpace(orderNr)
	trackingID = strings.TrimSpace(trackingID)
	trackingURL = strings.TrimSpace(trackingURL)
	carrier = strings.TrimSpace(carrier)

	if orderNr == "" {
		return ErrInvalidOrderNr
	}
	if trackingID == "" {
		return ErrInvalidTracking
	}

	o, err := s.client.Order.Query().
		Where(order.OrderNumber(orderNr)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return ErrOrderNotFound
		}
		return fmt.Errorf("fetch order: %w", err)
	}

	exists, err := s.client.Tracking.Query().
		Where(tracking.HasOrderWith(order.ID(o.ID))).
		Exist(ctx)
	if err != nil {
		return fmt.Errorf("check tracking exists: %w", err)
	}

	if !exists {
		_, err := s.client.Tracking.Create().
			SetOrder(o).
			SetTrackingID(trackingID).
			SetTrackingURL(trackingURL).
			SetCarrier(carrier).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("create tracking: %w", err)
		}
		return nil
	}

	t, err := s.client.Tracking.Query().
		Where(tracking.HasOrderWith(order.ID(o.ID))).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("load tracking: %w", err)
	}

	u := s.client.Tracking.UpdateOne(t).
		SetTrackingID(trackingID).
		SetTrackingURL(trackingURL)
	if carrier != "" {
		u.SetCarrier(carrier)
	}
	if err := u.Exec(ctx); err != nil {
		return fmt.Errorf("update tracking: %w", err)
	}
	return nil
}

func (s *TrackingService) Get(ctx context.Context, orderNr string) (*TrackingDTO, error) {
	orderNr = strings.TrimSpace(orderNr)
	if orderNr == "" {
		return nil, ErrInvalidOrderNr
	}

	o, err := s.client.Order.Query().
		Where(order.OrderNumber(orderNr)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("fetch order: %w", err)
	}

	t, err := s.client.Tracking.Query().
		Where(tracking.HasOrderWith(order.ID(o.ID))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return &TrackingDTO{OrderNr: o.OrderNumber}, nil
		}
		return nil, fmt.Errorf("fetch tracking: %w", err)
	}

	return &TrackingDTO{
		OrderNr:     o.OrderNumber,
		TrackingID:  t.TrackingID,
		TrackingURL: t.TrackingURL,
		Carrier:     t.Carrier,
	}, nil
}

func (s *TrackingService) Clear(ctx context.Context, orderNr string) error {
	orderNr = strings.TrimSpace(orderNr)
	if orderNr == "" {
		return ErrInvalidOrderNr
	}

	o, err := s.client.Order.Query().
		Where(order.OrderNumber(orderNr)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return ErrOrderNotFound
		}
		return fmt.Errorf("fetch order: %w", err)
	}

	n, err := s.client.Tracking.Delete().
		Where(tracking.HasOrderWith(order.ID(o.ID))).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("delte tracking: %w", err)
	}

	_ = n
	return nil
}
