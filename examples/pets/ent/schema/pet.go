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
			Annotations(elk.Annotation{
				// Include the name on the "pet:list" group.
				Groups: []string{"pet:list"},
			}),
	}
}

// Edges of the Pet.
func (Pet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("friends", Pet.Type),
		edge.From("owner", User.Type).
			Ref("pets").
			Unique().
			Annotations(elk.Annotation{
				// Include the owner on the "pet:list" group.
				Groups: []string{"pet:list"},
			}),
	}
}

// Annotations of the Pet.
func (Pet) Annotations() []schema.Annotation {
	return []schema.Annotation{
		elk.SchemaAnnotation{
			// Tell elk to use the "pet:list" group on list routes.
			ListGroups: []string{"pet:list"},
		},
	}
}
