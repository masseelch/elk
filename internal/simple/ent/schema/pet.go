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
		field.String("id").
			DefaultFunc(uuid.NewString),
		field.String("name").
			Annotations(elk.Groups("pet")),
		field.Int("age"),
	}
}

// Edges of the Pet.
func (Pet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("categories", Category.Type).
			Ref("pets"),
		edge.From("owner", Owner.Type).
			Ref("pets").
			Unique().
			Annotations(elk.Groups("pet:owner")),
		edge.To("friends", Pet.Type).
			Annotations(
				elk.Groups("pet"),
				elk.MaxDepth(3),
			),
	}
}

// Annotations of the Pet.
func (Pet) Annotations() []schema.Annotation {
	return []schema.Annotation{
		elk.ReadGroups("pet", "pet:owner", "owner"),
	}
}
