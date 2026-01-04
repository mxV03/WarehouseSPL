package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Tracking holds the schema definition for the Tracking entity.
type Tracking struct {
	ent.Schema
}

// Fields of the Tracking.
func (Tracking) Fields() []ent.Field {
	return []ent.Field{
		field.String("tracking_id").
			NotEmpty(),
		field.String("tracking_url").
			Optional().
			Default(""),
		field.String("carrier").
			Optional().
			Default(""),
		field.Time("created_at").
			Default(time.Now),
		field.Time("update_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

func (Tracking) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("order").
			Unique(),
	}
}

// Edges of the Tracking.
func (Tracking) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("order", Order.Type).
			Ref("tracking").
			Unique().
			Required(),
	}
}
