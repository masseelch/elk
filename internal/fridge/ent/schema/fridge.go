package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/masseelch/elk"
)

// Fridge holds the schema definition for the Fridge entity.
type Fridge struct {
	ent.Schema
}

// Fields of the Fridge.
func (Fridge) Fields() []ent.Field {
	return []ent.Field{
		field.String("title"),
	}
}

// Edges of the Fridge.
func (Fridge) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("compartments", Compartment.Type),
	}
}

// Annotations of the Fridge.
func (Fridge) Annotations() []schema.Annotation {
	return []schema.Annotation{
		elk.DeletePolicy(elk.Exclude),
	}
}
