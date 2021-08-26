package elk

import (
	"encoding/json"
	"entgo.io/ent/schema"
)

type (
	// SchemaAnnotation annotates an entity with metadata for templates.
	SchemaAnnotation struct {
		// CreateGroups holds the serializations groups to use on the creation handler.
		CreateGroups groups `json:",omitempty"`
		// ReadGroups holds the serializations groups to use on the read handler.
		ReadGroups groups `json:",omitempty"`
		// UpdateGroups holds the serializations groups to use on the update handler.
		UpdateGroups groups `json:",omitempty"`
		// ListGroups holds the serializations groups to use on the list handler.
		ListGroups groups `json:",omitempty"`
	}
	// Annotation annotates fields and edges with metadata for templates.
	Annotation struct {
		// Groups holds the serialization groups to use on this field / edge.
		Groups groups `json:",omitempty"`
		// MaxDepth tells the generator the maximum depth of this field when there is a cycle possible.
		MaxDepth uint
	}
)

// CreateGroups returns a create groups schema-annotation.
func CreateGroups(gs ...string) SchemaAnnotation {
	return SchemaAnnotation{CreateGroups: gs}
}

// ReadGroups returns a read groups schema-annotation.
func ReadGroups(gs ...string) SchemaAnnotation {
	return SchemaAnnotation{ReadGroups: gs}
}

// UpdateGroups returns an update groups schema-annotation.
func UpdateGroups(gs ...string) SchemaAnnotation {
	return SchemaAnnotation{UpdateGroups: gs}
}

// ListGroups returns a list groups schema-annotation.
func ListGroups(gs ...string) SchemaAnnotation {
	return SchemaAnnotation{ListGroups: gs}
}

// Groups returns a groups annotation.
func Groups(gs ...string) Annotation {
	return Annotation{Groups: gs}
}

// MaxDepth returns a max depth annotation.
func MaxDepth(d uint) Annotation {
	return Annotation{MaxDepth: d}
}

// Name implements ent.Annotation interface.
func (SchemaAnnotation) Name() string {
	return "ElkSchema"
}

// Merge implements ent.Merger interface.
func (a SchemaAnnotation) Merge(o schema.Annotation) schema.Annotation {
	var ant SchemaAnnotation
	switch o := o.(type) {
	case SchemaAnnotation:
		ant = o
	case *SchemaAnnotation:
		if o != nil {
			ant = *o
		}
	default:
		return a
	}
	if len(ant.CreateGroups) > 0 {
		a.CreateGroups = ant.CreateGroups
	}
	if len(ant.ReadGroups) > 0 {
		a.ReadGroups = ant.ReadGroups
	}
	if len(ant.UpdateGroups) > 0 {
		a.UpdateGroups = ant.UpdateGroups
	}
	if len(ant.ListGroups) > 0 {
		a.ListGroups = ant.ListGroups
	}
	return a
}

// Decode from ent.
func (a *SchemaAnnotation) Decode(o interface{}) error {
	buf, err := json.Marshal(o)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, a)
}

// Name implements ent.Annotation interface.
func (Annotation) Name() string {
	return "Elk"
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
	if len(ant.Groups) > 0 {
		a.Groups = ant.Groups
	}
	if ant.MaxDepth != 1 {
		a.MaxDepth = ant.MaxDepth
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

// EnsureDefaults ensures defaults are set.
func (a *Annotation) EnsureDefaults() {
	if a.MaxDepth == 0 {
		a.MaxDepth = 1
	}
}

var (
	_ schema.Annotation = (*Annotation)(nil)
	_ schema.Annotation = (*SchemaAnnotation)(nil)
)
