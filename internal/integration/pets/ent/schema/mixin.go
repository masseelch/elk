package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type ColorMixin struct {
	mixin.Schema
}

func (ColorMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("color").
			Values(
				"red", "orange", "yellow", "green", "blue", "indigo", "violet", "purple",
				"pink", "silver", "gold", "beige", "brown", "grey", "black", "white",
			),
	}
}

type MaterialMixin struct {
	mixin.Schema
}

func (MaterialMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("material").
			Values("leather", "plastic", "fabric"),
	}
}
