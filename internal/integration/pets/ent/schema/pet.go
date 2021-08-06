package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/masseelch/elk"
)

// Pet holds the schema definition for the Pet entity.
type Pet struct {
	ent.Schema
}

// Fields of the Pet.
func (Pet) Fields() []ent.Field {
	return []ent.Field{
		field.Int("height").
			Positive(),
		field.Float("weight").
			Positive(),
		field.Bool("castrated"),
		field.String("name").
			Unique(),
		field.Time("birthday"),
		field.JSON("nicknames", []string{}).
			Optional(),
		field.Enum("sex").
			Values("male", "female"),
		field.UUID("chip", uuid.UUID{}).
			Default(uuid.New),
	}
}

// Edges of the Pet.
func (Pet) Edges() []ent.Edge {
	return []ent.Edge{
		// O20 - Two Types
		edge.To("badge", Badge.Type).
			Unique().
			Annotations(
				elk.Groups("pet:read", "pet:list"),
			),
		// O20 - Same Types
		edge.To("mentor", Pet.Type).
			Unique().
			From("protege").
			Unique().
			Annotations(
				elk.Groups("pet:read"),
			),
		// O20 - Bidirectional
		edge.To("spouse", Pet.Type).
			Unique().
			Annotations(
				elk.Groups("pet:read"),
			),
		// O2M - Two Types
		edge.To("toys", Toy.Type).
			Annotations(
				elk.Groups("pet:read"),
			),
		// O2M - Same Types
		edge.To("children", Pet.Type).
			From("parent").
			Unique().
			Annotations(
				elk.Groups("pet:read"),
			),
		// M2M - Two Types
		edge.To("play_groups", PlayGroup.Type).
			Annotations(
				elk.Groups("pet:read"),
			),
		// M2M - Same Type - no idea
		// M2M - Bidirectional
		edge.To("friends", Pet.Type).
			Annotations(
				elk.Groups("pet:read"),
				elk.MaxDepth(2),
			),
	}
}

// Annotations of the Pet.
func (Pet) Annotations() []schema.Annotation {
	return []schema.Annotation{
		elk.ListGroups("pet:list"),
		elk.ReadGroups("pet:read"),
	}
}
