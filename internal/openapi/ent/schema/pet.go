package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/masseelch/elk"
	"github.com/masseelch/elk/spec"
)

// Pet holds the schema definition for the Pet entity.
type Pet struct {
	ent.Schema
}

// Fields of the Pet.
func (Pet) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Annotations(
				elk.Groups("pet"),
				elk.Example("Kuro"),
			),
		field.JSON("nicknames", []string{}).
			Optional().
			Annotations(elk.Groups("pet:read")),
		field.Int("age").
			Optional().
			Annotations(
				elk.Example(1),
			),
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
		elk.SchemaSecurity(spec.Security{{"apiKeySample": {}}}),
	}
}
