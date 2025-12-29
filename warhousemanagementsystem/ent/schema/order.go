package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Order holds the schema definition for the Order entity.
type Order struct {
	ent.Schema
}

// Fields of the Order.
func (Order) Fields() []ent.Field {
	return []ent.Field{
		field.String("order_number").
			Unique().
			NotEmpty(),
		field.String("type").
			NotEmpty(), // "INBOUND" or "OUTBOUND"
		field.String("status").
			NotEmpty(). // "DRAFT", "POSTED", "CANCELLED"
			Default("DRAFT"),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the Order.
func (Order) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("lines", OrderLine.Type),
	}
}
