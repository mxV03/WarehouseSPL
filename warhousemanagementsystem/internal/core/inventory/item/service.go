package item

import (
	"context"
	"fmt"
	"strings"

	"github.com/mxV03/warhousemanagementsystem/ent"
	"github.com/mxV03/warhousemanagementsystem/ent/item"
)

var (
	ErrItemNotFound = fmt.Errorf("item not found")
	ErrInvalidSKU   = fmt.Errorf("invalid SKU")
	ErrInvalidName  = fmt.Errorf("invalid name")
	ErrItemExists   = fmt.Errorf("item already exists")
)

type ItemService struct {
	client *ent.Client
}

func NewItemService(client *ent.Client) *ItemService {
	return &ItemService{client: client}
}

type ItemDTO struct {
	ID          int
	SKU         string
	Name        string
	Description string
}

func (s *ItemService) CreateItem(ctx context.Context, sku, name, description string) (*ItemDTO, error) {
	sku = strings.TrimSpace(sku)
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)

	if sku == "" {
		return nil, ErrInvalidSKU
	}
	if name == "" {
		return nil, ErrInvalidName
	}

	exists, err := s.client.Item.Query().Where(item.SKU(sku)).Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("checking item existence: %w", err)
	}

	if exists {
		return nil, ErrItemExists
	}

	itm, err := s.client.Item.Create().
		SetSKU(sku).
		SetName(name).
		SetDescription(description).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating item: %w", err)
	}

	return &ItemDTO{
		ID:          itm.ID,
		SKU:         itm.SKU,
		Name:        itm.Name,
		Description: itm.Description,
	}, nil
}

func (s *ItemService) GetItemByID(ctx context.Context, id int) (*ItemDTO, error) {
	itm, err := s.client.Item.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrItemNotFound
		}
		return nil, fmt.Errorf("getting item by ID: %w", err)
	}
	return &ItemDTO{
		ID:          itm.ID,
		SKU:         itm.SKU,
		Name:        itm.Name,
		Description: itm.Description,
	}, nil
}

func (s *ItemService) GetItemBySKU(ctx context.Context, sku string) (*ItemDTO, error) {
	sku = strings.TrimSpace(sku)
	if sku == "" {
		return nil, ErrInvalidSKU
	}
	itm, err := s.client.Item.Query().Where(item.SKU(sku)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrItemNotFound
		}
		return nil, fmt.Errorf("getting item by SKU: %w", err)
	}

	return &ItemDTO{
		ID:          itm.ID,
		SKU:         itm.SKU,
		Name:        itm.Name,
		Description: itm.Description,
	}, nil
}

func (s *ItemService) ListItems(ctx context.Context, limit int) ([]*ItemDTO, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	items, err := s.client.Item.Query().Order(ent.Asc(item.FieldSKU)).Limit(limit).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing items: %w", err)
	}

	out := make([]*ItemDTO, 0, len(items))
	for _, itm := range items {
		out = append(out, &ItemDTO{
			ID:          itm.ID,
			SKU:         itm.SKU,
			Name:        itm.Name,
			Description: itm.Description,
		})
	}
	return out, nil

}

func (s *ItemService) DeleteItemBySKU(ctx context.Context, sku string) error {
	sku = strings.TrimSpace(sku)
	if sku == "" {
		return ErrInvalidSKU
	}

	deleted, err := s.client.Item.Delete().Where(item.SKU(sku)).Exec(ctx)
	if err != nil {
		return fmt.Errorf("deleting item by SKU: %w", err)
	}

	if deleted == 0 {
		return ErrItemNotFound
	}
	return nil
}
