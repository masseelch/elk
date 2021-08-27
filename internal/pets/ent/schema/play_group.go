package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// PlayGroup holds the schema definition for the PlayGroup entity.
type PlayGroup struct {
	ent.Schema
}

// Fields of the PlayGroup.
func (PlayGroup) Fields() []ent.Field {
	return []ent.Field{
		field.String("title"),
		field.Text("description").
			Optional(),
		field.Enum("weekday").
			Values("Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"),
	}
}

// Edges of the PlayGroup.
func (PlayGroup) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("participants", Pet.Type).
			Ref("play_groups"),
	}
}
