package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// OrderLine holds the schema definition for the OrderLine entity.
type OrderLine struct {
	ent.Schema
}

// Fields of the OrderLine.
func (OrderLine) Fields() []ent.Field {
	return []ent.Field{
		field.Int("quantity").
			Positive(),
	}
}

// Edges of the OrderLine.
func (OrderLine) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("order", Order.Type).
			Ref("lines").
			Unique().
			Required(),

		edge.From("item", Item.Type).
			Ref("order_lines").
			Unique().
			Required(),

		edge.From("location", Location.Type).
			Ref("order_lines").
			Unique().
			Required(),
	}
}
