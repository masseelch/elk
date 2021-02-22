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
		field.String("name").StructTag(`groups:"pet:list"`),
		field.Int("age").Optional().Nillable().StructTag(`groups:"pet:list"`),
		field.Uint32("color").GoType(Color(0)).StructTag(`groups:"pet:list"`),
		field.Strings("string_list").StructTag(`groups:"pet:list"`),
	}
}

// Edges of the Pet.
func (Pet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", Owner.Type).
			Ref("pets").
			Unique().
			StructTag(`json:"owner" groups:"pet:list"`),
	}
}

// Annotations of the Pet.
func (Pet) Annotations() []schema.Annotation {
	return []schema.Annotation{
		edge.Annotation{StructTag: `json:"edges" groups:"pet:list"`},
		elk.HandlerAnnotation{ListGroups: []string{"pet:list", "owner:list"}},
	}
}

type Color uint32
