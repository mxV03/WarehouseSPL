package orders

import (
	"context"
	"fmt"

	"github.com/mxV03/warhousemanagementsystem/ent"
	"github.com/mxV03/warhousemanagementsystem/ent/item"
	"github.com/mxV03/warhousemanagementsystem/ent/location"
	"github.com/mxV03/warhousemanagementsystem/ent/order"
	"github.com/mxV03/warhousemanagementsystem/ent/orderline"
	"github.com/mxV03/warhousemanagementsystem/internal/core/inventory/stock"
)

var (
	ErrOrderNotFound    = fmt.Errorf("order not found")
	ErrOrderExists      = fmt.Errorf("order already exists")
	ErrInvalidOrderNo   = fmt.Errorf("invalid order number")
	ErrInvalidOrderType = fmt.Errorf("invalid order type")
	ErrInvalidStatus    = fmt.Errorf("invalid order status transaction")
	ErrNoLines          = fmt.Errorf("order has no lines")
	ErrInvalidQuantity  = fmt.Errorf("invalid quantity specified")
)

type OrderType string

const (
	OrderTypeInbound  OrderType = "INBOUND"
	OrderTypeOutbound OrderType = "OUTBOUND"
)

type OrderStatus string

const (
	OrderStatusDraft     OrderStatus = "DRAFT"
	OrderStatusPosted    OrderStatus = "POSTED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

type OrderDTO struct {
	Id        int
	Number    string
	Type      string
	Status    string
	CreatedAt string
}

type OrderLineDTO struct {
	Id           int
	OrderNumber  string
	SKU          string
	LocationCode string
	Quantity     int
}

type OrderService struct {
	client *ent.Client
}

func NewOrderService(client *ent.Client) *OrderService {
	return &OrderService{client: client}
}

func (s *OrderService) CreateInboundOrder(ctx context.Context, number string) (*ent.Order, error) {
	return s.create(ctx, number, string(OrderTypeInbound))
}

func (s *OrderService) CreateOutboundOrder(ctx context.Context, number string) (*ent.Order, error) {
	return s.create(ctx, number, string(OrderTypeOutbound))
}

func (s *OrderService) create(ctx context.Context, number string, orderType string) (*ent.Order, error) {
	exists, err := s.client.Order.Query().
		Where(order.OrderNumber(number)).
		Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("checking order existence: %w", err)
	}
	if exists {
		return nil, ErrOrderExists
	}

	newOrder := s.client.Order.Create().
		SetOrderNumber(number).
		SetType(string(orderType)).
		SetStatus(string(OrderStatusDraft))

	createdOrder, err := newOrder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating order: %w", err)
	}
	return createdOrder, nil
}

func (s *OrderService) AddLine(ctx context.Context, orderNumber, sku, locationCode string, quantity int) (*ent.OrderLine, error) {
	if orderNumber == "" {
		return nil, ErrInvalidOrderNo
	}
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	orderEntity, err := s.client.Order.Query().
		Where(order.OrderNumber(orderNumber)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("fetching order: %w", err)
	}
	if orderEntity.Status != string(OrderStatusDraft) {
		return nil, ErrInvalidStatus
	}

	itm, err := s.client.Item.Query().
		Where(item.SKU(sku)).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching item: %w", err)
	}

	loc, err := s.client.Location.Query().
		Where(location.Code(locationCode)).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching location: %w", err)
	}

	line := s.client.OrderLine.Create().
		SetOrder(orderEntity).
		SetItem(itm).
		SetLocation(loc).
		SetQuantity(quantity)
	createdLine, err := line.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating order line: %w", err)
	}
	return createdLine, nil
}

func (s *OrderService) PostOrder(ctx context.Context, number string) error {
	if number == "" {
		return ErrInvalidOrderNo
	}

	tx, err := s.client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback()

	orderEntity, err := tx.Order.Query().
		Where(order.OrderNumber(number)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return ErrOrderNotFound
		}
		return fmt.Errorf("fetching order: %w", err)
	}
	if orderEntity.Status != string(OrderStatusDraft) {
		return ErrInvalidStatus
	}

	lines, err := tx.OrderLine.Query().
		Where(orderline.HasOrderWith(order.OrderNumber(number))).
		WithItem().
		WithLocation().
		All(ctx)

	if err != nil {
		return fmt.Errorf("fetching order lines: %w", err)
	}
	if len(lines) == 0 {
		return ErrNoLines
	}

	stockSvc := stock.NewStockService(tx.Client())
	for _, line := range lines {
		sku := line.Edges.Item.SKU
		locCode := line.Edges.Location.Code
		qty := line.Quantity

		ref := "ORDER-" + number

		if orderEntity.Type == string(OrderTypeInbound) {
			if err := stockSvc.IN(ctx, sku, locCode, qty, ref); err != nil {
				return err
			}
		} else if orderEntity.Type == string(OrderTypeOutbound) {
			if err := stockSvc.OUT(ctx, sku, locCode, qty, ref); err != nil {
				return err
			}
		} else {
			return ErrInvalidOrderType
		}
	}

	_, err = tx.Order.UpdateOneID(orderEntity.ID).
		SetStatus(string(OrderStatusPosted)).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("updating order status: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}

func (s *OrderService) CancelOrder(ctx context.Context, orderNumber string) error {
	if orderNumber == "" {
		return ErrInvalidOrderNo
	}

	o, err := s.client.Order.Query().Where(order.OrderNumber(orderNumber)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return ErrOrderNotFound
		}
		return fmt.Errorf("fetching order: %w", err)
	}
	if o.Status != string(OrderStatusDraft) {
		return ErrInvalidStatus
	}

	_, err = s.client.Order.UpdateOneID(o.ID).
		SetStatus(string(OrderStatusCancelled)).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("cancel order: %w", err)
	}
	return nil
}
