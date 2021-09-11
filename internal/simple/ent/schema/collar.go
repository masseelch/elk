package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Collar holds the schema definition for the Collar entity.
type Collar struct {
	ent.Schema
}

// Fields of the Collar.
func (Collar) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("color").
			Values("green", "red", "blue"),
	}
}

// Edges of the Collar.
func (Collar) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("pet", Pet.Type).
			Ref("collar").
			Unique(),
	}
}
