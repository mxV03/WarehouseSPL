package schema

import "entgo.io/ent"

// Item holds the schema definition for the Item entity.
type Item struct {
	ent.Schema
}

// Fields of the Item.
func (Item) Fields() []ent.Field {
	return nil
}

// Edges of the Item.
func (Item) Edges() []ent.Edge {
	return nil
}
