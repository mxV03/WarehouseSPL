package schema

import (
	"entgo.io/ent"
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
	return nil
}
