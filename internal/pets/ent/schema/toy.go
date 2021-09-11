package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Toy holds the schema definition for the Toy entity.
type Toy struct {
	ent.Schema
}

// Fields of the Toy.
func (Toy) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("title"),
	}
}

// Mixin of the Toy.
func (Toy) Mixin() []ent.Mixin {
	return []ent.Mixin{
		ColorMixin{},
		MaterialMixin{},
	}
}

// Edges of the Toy.
func (Toy) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", Pet.Type).
			Ref("toys").
			Unique(),
	}
}
