package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// AuditEvent holds the schema definition for the AuditEvent entity.
type AuditEvent struct {
	ent.Schema
}

// Fields of the AuditEvent.
func (AuditEvent) Fields() []ent.Field {
	return []ent.Field{
		field.Time("ts").
			Default(time.Now),
		field.String("action").
			NotEmpty(),
		field.String("entity").
			NotEmpty(),
		field.String("entity_ref").
			Optional().
			Default(""),
		field.String("actor").
			Optional().
			Default("system"),
		field.String("details").
			Optional().
			Default(""),
	}
}

func (AuditEvent) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("ts"),
		index.Fields("action"),
		index.Fields("entity"),
		index.Fields("actor"),
	}
}

// Edges of the AuditEvent.
func (AuditEvent) Edges() []ent.Edge {
	return nil
}
