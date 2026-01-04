package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Warehouse holds the schema definition for the Warehouse entity.
type Warehouse struct {
	ent.Schema
}

// Fields of the Warehouse.
func (Warehouse) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").NotEmpty(),
		field.String("name").Optional().Default(""),
	}
}

func (Warehouse) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("code").Unique(),
	}
}

// Edges of the Warehouse.
func (Warehouse) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("warehouse_locations", WarehouseLocation.Type),
	}
}
