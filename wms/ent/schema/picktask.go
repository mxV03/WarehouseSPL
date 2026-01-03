package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// PickTask holds the schema definition for the PickTask entity.
type PickTask struct {
	ent.Schema
}

// Fields of the PickTask.
func (PickTask) Fields() []ent.Field {
	return []ent.Field{
		field.Int("quantity").Positive(),
		field.String("status").Default("OPEN"),
		field.Time("picked_at").Optional().Nillable(),
		field.Int("bin_id").Optional().Nillable(),
	}
}

func (PickTask) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("picklist", "order_line").Unique(),
	}
}

// Edges of the PickTask.
func (PickTask) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("picklist", PickList.Type).
			Ref("tasks").
			Unique().
			Required(),

		edge.From("order_line", OrderLine.Type).
			Ref("pick_tasks").
			Unique().
			Required(),
	}
}
