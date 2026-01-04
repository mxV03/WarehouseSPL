package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/index"
)

// WarehouseLocation holds the schema definition for the WarehouseLocation entity.
type WarehouseLocation struct {
	ent.Schema
}

// Fields of the WarehouseLocation.
func (WarehouseLocation) Fields() []ent.Field {
	return nil
}

func (WarehouseLocation) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("location").Unique(),

		index.Edges("warehouse", "location").Unique(),
	}
}

// Edges of the WarehouseLocation.
func (WarehouseLocation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("warehouse", Warehouse.Type).
			Ref("warehouse_locations").
			Unique().
			Required(),
		edge.From("location", Location.Type).
			Ref("warehouse_link").
			Unique().
			Required(),
	}
}
