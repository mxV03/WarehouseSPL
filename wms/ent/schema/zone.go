package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Zone holds the schema definition for the Zone entity.
type Zone struct {
	ent.Schema
}

// Fields of the Zone.
func (Zone) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").
			NotEmpty(),
		field.String("name").
			Optional(),
	}
}

func (Zone) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("code").
			Edges("location").
			Unique(),
	}
}

// Edges of the Zone.
func (Zone) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("location", Location.Type).
			Ref("zones").
			Unique().
			Required(),

		edge.To("bins", Bin.Type),
	}
}
