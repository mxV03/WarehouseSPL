package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Bin holds the schema definition for the Bin entity.
type Bin struct {
	ent.Schema
}

// Fields of the Bin.
func (Bin) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").
			NotEmpty(),
		field.String("name").Optional(),
	}
}

func (Bin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("code").
			Edges("location").
			Unique(),
	}
}

// Edges of the Bin.
func (Bin) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("location", Location.Type).
			Ref("bins").
			Unique().
			Required(),

		edge.From("zone", Zone.Type).
			Ref("bins").
			Unique().
			Required(),

		edge.To("items", Item.Type),
	}
}
