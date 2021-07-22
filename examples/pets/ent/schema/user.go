package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/masseelch/elk"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("age").
			Annotations(elk.Annotation{
				// Include the age only on the "user:read" group.
				Groups: []string{"user:read"},
			}),
		field.String("name").
			Annotations(elk.Annotation{
				// No numbers allowed in name and it has to be at least 3 chars long.
				CreateValidation: "alpha,min=3",
				UpdateValidation: "alpha,min=3",
			}),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("pets", Pet.Type).
			Annotations(elk.Annotation{
				// We'd like to see the pets rendered if the "user:read" group is given.
				Groups: []string{"user:read"},
			}),
		edge.To("friends", User.Type).
			Annotations(elk.Annotation{
				// We'd like to see the friends rendered if the "user:read" group is given.
				Groups: []string{"user:read"},
				// Give us 2 levels of friends.
				MaxDepth: 2,
			}),
		edge.From("groups", Group.Type).
			Ref("users").
			Annotations(elk.Annotation{
				// What groups does this user belong to (again serialize only on "user:read").
				Groups: []string{"user:read"},
			}),
		edge.From("manage", Group.Type).
			Ref("admin").
			Annotations(elk.Annotation{
				// What groups does this user manage?
				Groups: []string{"user:read"},
			}),
	}
}

// Annotations of the User.
func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		elk.SchemaAnnotation{
			// Tell elk to use the "user:read" group on read routes.
			ReadGroups: []string{"user:read"},
		},
	}
}
