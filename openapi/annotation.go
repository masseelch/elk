package openapi

import (
	"encoding/json"
	"entgo.io/ent/schema"
)

// Annotation can be used on either schemas or field / edges to add information
// about some OpenAPI specs like tags or field examples.
type Annotation struct {
	// OpenAPI spec tags on routes.
	Tags []string `json:",omitempty"`
	// OpenAPI spec example value on schema fields.
	Example interface{} `json:",omitempty"`
}

// Tags returns a tags annotation.
func Tags(tags ...string) Annotation {
	return Annotation{Tags: tags}
}

// Example returns an example annotation.
func Example(v interface{}) Annotation {
	return Annotation{Example: v}
}

// Name implements ent.Annotation interface.
func (Annotation) Name() string {
	return "OpenAPI"
}

// Merge implements ent.Merger interface.
func (a Annotation) Merge(o schema.Annotation) schema.Annotation {
	var ant Annotation
	switch o := o.(type) {
	case Annotation:
		ant = o
	case *Annotation:
		if o != nil {
			ant = *o
		}
	default:
		return a
	}
	if len(ant.Tags) > 0 {
		a.Tags = ant.Tags
	}
	if ant.Example != nil {
		a.Example = ant.Example
	}
	return a
}

// Decode from ent.
func (a *Annotation) Decode(o interface{}) error {
	buf, err := json.Marshal(o)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, a)
}

var _ schema.Annotation = (*Annotation)(nil)
