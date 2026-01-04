package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Location holds the schema definition for the Location entity.
type Location struct {
	ent.Schema
}

// Fields of the Location.
func (Location) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").
			Unique().
			NotEmpty(),
		field.String("name").
			NotEmpty(),
	}
}

// Edges of the Location.
func (Location) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("movements", StockMovement.Type),
		edge.To("order_lines", OrderLine.Type),

		edge.To("zones", Zone.Type),
		edge.To("bins", Bin.Type),
		edge.To("warehouse_link", WarehouseLocation.Type).Unique(),
	}
}
