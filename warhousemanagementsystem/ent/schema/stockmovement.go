package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// StockMovement holds the schema definition for the StockMovement entity.
type StockMovement struct {
	ent.Schema
}

// Fields of the StockMovement.
func (StockMovement) Fields() []ent.Field {
	return []ent.Field{
		field.String("type").
			NotEmpty(), // "IN" or "OUT"
		field.Int("quantity").
			Positive(),
		field.Time("created_at").
			Default(time.Now),
		field.String("reference").
			Optional(),
	}
}

// Edges of the StockMovement.
func (StockMovement) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("item", Item.Type).
			Ref("movements").
			Unique().
			Required(),

		edge.From("location", Location.Type).
			Ref("movements").
			Unique().
			Required(),
	}
}
