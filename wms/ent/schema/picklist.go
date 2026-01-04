package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// PickList holds the schema definition for the PickList entity.
type PickList struct {
	ent.Schema
}

// Fields of the PickList.
func (PickList) Fields() []ent.Field {
	return []ent.Field{
		field.String("status").
			Default("CREATED"),
		field.Time("created_at").
			Default(time.Now),
		field.Time("started_at").
			Optional().
			Nillable(),
		field.Time("done_at").
			Optional().
			Nillable(),
	}
}

func (PickList) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("order").Unique(),
	}
}

// Edges of the PickList.
func (PickList) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("order", Order.Type).
			Ref("picklist").
			Unique().
			Required(),

		edge.To("tasks", PickTask.Type),
	}
}
