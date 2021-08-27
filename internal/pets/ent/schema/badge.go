package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
)

// Badge holds the schema definition for the Badge entity.
type Badge struct {
	ent.Schema
}

// Mixin of the Badge.
func (Badge) Mixin() []ent.Mixin {
	return []ent.Mixin{
		ColorMixin{},
		MaterialMixin{},
	}
}

// Edges of the Badge.
func (Badge) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("wearer", Pet.Type).
			Ref("badge").
			Unique(),
	}
}
