package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/masseelch/elk"
)

// Owner holds the schema definition for the Owner entity.
type Owner struct {
	ent.Schema
}

// Fields of the Owner.
func (Owner) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Annotations(
				elk.Annotation{
					Groups: []string{"owner"},
				},
			),
		field.String("age").
			Annotations(
				elk.Annotation{
					Groups: []string{"owner"},
				},
			),
	}
}

// Edges of the Owner.
func (Owner) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("pets", Pet.Type),
		edge.To("friends", Owner.Type).
			Annotations(
				elk.Annotation{
					Groups: []string{"owner"},
				},
			),
	}
}
