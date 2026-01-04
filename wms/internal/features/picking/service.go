//go:build picking

package picking

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/mxV03/wms/ent"
	"github.com/mxV03/wms/ent/bin"
	"github.com/mxV03/wms/ent/item"
	"github.com/mxV03/wms/ent/location"
	"github.com/mxV03/wms/ent/order"
	"github.com/mxV03/wms/ent/picklist"
	"github.com/mxV03/wms/ent/picktask"
)

var (
	ErrInvalidOrderNr   = fmt.Errorf("invalid order number")
	ErrPickListExists   = fmt.Errorf("picklist already exists for order")
	ErrPickListNotFound = fmt.Errorf("picklist not found")
	ErrTaskNotFound     = fmt.Errorf("pick task not found")
	ErrInvalidStatus    = fmt.Errorf("invalid status transition")
)

type PickingService struct {
	client *ent.Client
}

func NewPickingService(client *ent.Client) *PickingService {
	return &PickingService{
		client: client,
	}
}

type TaskDTO struct {
	ID       int
	SKU      string
	ItemName string
	Location string
	Bin      string
	Quantity int
	Status   string
}

type PickListDTO struct {
	ID        int
	OrderNr   string
	Status    string
	CreatedAt time.Time
	StartedAt *time.Time
	DoneAt    *time.Time
	Tasks     []TaskDTO
}

func (s *PickingService) CreatePickList(ctx context.Context, orderNr string) (*ent.PickList, error) {
	orderNr = strings.TrimSpace(orderNr)
	if orderNr == "" {
		return nil, ErrInvalidOrderNr
	}

	o, err := s.client.Order.Query().
		Where(order.OrderNumber(orderNr)).
		WithLines(func(q *ent.OrderLineQuery) {
			q.WithItem().WithLocation()
		}).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("fetch order: %w", err)
	}

	exists, err := s.client.PickList.Query().
		Where(picklist.HasOrderWith(order.OrderNumber(orderNr))).
		Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("check picklist exists: %w", err)
	}
	if exists {
		return nil, ErrPickListExists
	}

	pl, err := s.client.PickList.Create().
		SetOrder(o).
		SetStatus("CREATED").
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, ErrPickListExists
		}
		return nil, fmt.Errorf("create picklist: %w", err)
	}

	for _, ol := range o.Edges.Lines {
		taskCreate := s.client.PickTask.Create().
			SetPicklist(pl).
			SetOrderLine(ol).
			SetQuantity(ol.Quantity).
			SetStatus("OPEN")

		binID, _ := s.getBinIDForLine(ctx, ol)
		if binID != nil {
			taskCreate.SetBinID(*binID)
		}

		if _, err := taskCreate.Save(ctx); err != nil {
			return nil, fmt.Errorf("create pick task: %w", err)
		}
	}

	return pl, nil
}

func (s *PickingService) getBinIDForLine(ctx context.Context, ol *ent.OrderLine) (*int, error) {
	if ol == nil || ol.Edges.Item == nil {
		return nil, nil
	}

	q := s.client.Bin.Query().
		Where(bin.HasItemsWith(item.ID(ol.Edges.Item.ID)))

	if ol.Edges.Location != nil {
		q = q.Where(bin.HasLocationWith(location.ID(ol.Edges.Location.ID)))
	}

	b, err := q.First(ctx)
	if err == nil {
		id := b.ID
		return &id, nil
	}
	return nil, nil
}

func (s *PickingService) StartPickList(ctx context.Context, pickListID int) error {
	pl, err := s.client.PickList.Get(ctx, pickListID)
	if err != nil {
		if ent.IsNotFound(err) {
			return ErrPickListNotFound
		}
		return fmt.Errorf("fetch picklist: %w", err)
	}
	if pl.Status != "CREATED" {
		return ErrInvalidStatus
	}
	now := time.Now()
	return s.client.PickList.UpdateOne(pl).
		SetStatus("IN_PROGRESS").
		SetStartedAt(now).
		Exec(ctx)
}

func (s *PickingService) MarkTaskPicked(ctx context.Context, taskID int) error {
	t, err := s.client.PickTask.Get(ctx, taskID)
	if err != nil {
		if ent.IsNotFound(err) {
			return ErrTaskNotFound
		}
		return fmt.Errorf("fetch task: %w", err)
	}
	if t.Status != "OPEN" {
		return ErrInvalidStatus
	}
	now := time.Now()
	return s.client.PickTask.UpdateOne(t).
		SetStatus("PICKED").
		SetPickedAt(now).
		Exec(ctx)
}

func (s *PickingService) DonePickList(ctx context.Context, pickListID int) error {
	pl, err := s.client.PickList.Get(ctx, pickListID)
	if err != nil {
		if ent.IsNotFound(err) {
			return ErrPickListNotFound
		}
		return fmt.Errorf("fetch picklist: %w", err)
	}
	if pl.Status != "IN_PROGRESS" {
		return ErrInvalidStatus
	}

	openCount, err := s.client.PickTask.Query().
		Where(picktask.HasPicklistWith(picklist.ID(pickListID)), picktask.StatusEQ("OPEN")).
		Count(ctx)
	if err != nil {
		return fmt.Errorf("count open tasks: %w", err)
	}
	if openCount > 0 {
		return fmt.Errorf("cannot finish picklist: %d task(s) still OPEN", openCount)
	}

	now := time.Now()
	return s.client.PickList.UpdateOneID(pickListID).
		SetStatus("DONE").
		SetDoneAt(now).
		Exec(ctx)
}

func (s *PickingService) ShowPickList(ctx context.Context, pickListID int) (*PickListDTO, error) {
	pl, err := s.client.PickList.Query().
		Where(picklist.ID(pickListID)).
		WithOrder().
		WithTasks(func(tq *ent.PickTaskQuery) {
			tq.WithOrderLine(func(olq *ent.OrderLineQuery) {
				olq.WithItem().WithLocation()
			})
		}).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrPickListNotFound
		}
		return nil, fmt.Errorf("fetch picklist: %w", err)
	}

	dto := &PickListDTO{
		ID:        pl.ID,
		Status:    pl.Status,
		CreatedAt: pl.CreatedAt,
		StartedAt: pl.StartedAt,
		DoneAt:    pl.DoneAt,
	}
	if pl.Edges.Order != nil {
		dto.OrderNr = pl.Edges.Order.OrderNumber
	}

	binIDs := make([]int, 0)
	seen := map[int]struct{}{}
	for _, t := range pl.Edges.Tasks {
		if t.BinID != nil {
			if _, ok := seen[*t.BinID]; !ok {
				seen[*t.BinID] = struct{}{}
				binIDs = append(binIDs, *t.BinID)
			}
		}
	}
	sort.Ints(binIDs)

	binCodeByID := map[int]string{}
	if len(binIDs) > 0 {
		bins, err := s.client.Bin.Query().
			Where(bin.IDIn(binIDs...)).
			All(ctx)
		if err == nil {
			for _, b := range bins {
				binCodeByID[b.ID] = b.Code
			}
		}
	}

	for _, t := range pl.Edges.Tasks {
		td := TaskDTO{
			ID:       t.ID,
			Quantity: t.Quantity,
			Status:   t.Status,
			Bin:      "-",
		}

		if t.BinID != nil {
			if code, ok := binCodeByID[*t.BinID]; ok && strings.TrimSpace(code) != "" {
				td.Bin = code
			}
		}

		if t.Edges.OrderLine != nil {
			ol := t.Edges.OrderLine
			if ol.Edges.Item != nil {
				td.SKU = ol.Edges.Item.SKU
				td.ItemName = ol.Edges.Item.Name
			}
			if ol.Edges.Location != nil {
				td.Location = ol.Edges.Location.Code
			}
		}
		dto.Tasks = append(dto.Tasks, td)
	}
	return dto, nil
}
