package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Item holds the schema definition for the Item entity.
type Item struct {
	ent.Schema
}

// Fields of the Item.
func (Item) Fields() []ent.Field {
	return []ent.Field{
		field.String("SKU").
			Unique().
			NotEmpty(),
		field.String("name").
			NotEmpty(),
		field.String("description").
			Optional(),
	}
}

// Edges of the Item.
func (Item) Edges() []ent.Edge {
	return nil
}
