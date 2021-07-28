package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/masseelch/elk"
)

// Pet holds the schema definition for the Pet entity.
type Pet struct {
	ent.Schema
}

// Fields of the Pet.
func (Pet) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			MinLen(3).
			Annotations(
				elk.Annotation{
					Groups: []string{"pet"},
				},
			),
		field.Int("age").
			Positive(),
	}
}

// Edges of the Pet.
func (Pet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("category", Category.Type).
			Ref("pets"),
		edge.From("owner", Owner.Type).
			Ref("pets").
			Unique().
			Annotations(
				elk.Annotation{
					Groups: []string{"pet:owner"},
				},
			),
		edge.To("friends", Pet.Type).
			Annotations(
				elk.Annotation{
					Groups:   []string{"pet"},
					MaxDepth: 3,
				},
			),
	}
}

// Annotations of the Pet.
func (Pet) Annotations() []schema.Annotation {
	return []schema.Annotation{
		elk.SchemaAnnotation{
			ReadGroups: []string{"pet", "pet:owner", "owner"},
		},
	}
}
