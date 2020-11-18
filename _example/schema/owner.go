package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema"
	"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
	"github.com/masseelch/elk"
)

// Owner holds the schema definition for the Owner entity.
type Owner struct {
	ent.Schema
}

// Fields of the Owner.
func (Owner) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").StructTag(`groups:"owner:read"`),
	}
}

// Edges of the Owner.
func (Owner) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("pets", Pet.Type).
			StructTag(`json:"pets" groups:"owner:read"`).
			Annotations(
				elk.EdgeAnnotation{
					DefaultOrder: []elk.Order{
						{Order: "asc", Field: "name"},
						{Order: "desc", Field: "id"},
					},
				},
			),
	}
}

// Annotations of the Owner.
func (Owner) Annotations() []schema.Annotation {
	return []schema.Annotation{
		edge.Annotation{
			StructTag: `json:"edges" groups:"owner:read"`,
		},
		elk.HandlerAnnotation{
			ReadGroups: []string{
				"owner:read",
				"pet:list",
			},
			DefaultListOrder: []elk.Order{{Order: "desc", Field: "name"}},
		},
	}
}
