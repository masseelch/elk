package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema"
	"github.com/facebook/ent/schema/field"
	"github.com/masseelch/elk"
)

// SkipGenerationModel holds the schema definition for the SkipGenerationModel entity.
type SkipGenerationModel struct {
	ent.Schema
}

// Fields of the SkipGenerationModel.
func (SkipGenerationModel) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").StructTag(`groups:"SkipGenerationModel:list"`),
	}
}

// Annotations of the SkipGenerationModel.
func (SkipGenerationModel) Annotations() []schema.Annotation {
	return []schema.Annotation{
		elk.HandlerAnnotation{Skip: true},
	}
}