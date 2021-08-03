package elk

import (
	"encoding/json"
	"entgo.io/ent/schema"
)

type (
	// SchemaAnnotation annotates an entity with metadata for templates.
	SchemaAnnotation struct {
		// CreateGroups holds the serializations groups to use on the create handler.
		CreateGroups groups `json:"CreateGroups,omitempty"`
		// ReadGroups holds the serializations groups to use on the read handler.
		ReadGroups groups `json:"ReadGroups,omitempty"`
		// UpdateGroups holds the serializations groups to use on the update handler.
		UpdateGroups groups `json:"UpdateGroups,omitempty"`
		// ListGroups holds the serializations groups to use on the list handler.
		ListGroups groups `json:"ListGroups,omitempty"`
	}
	// Annotation annotates fields and edges with metadata for templates.
	Annotation struct {
		// Groups holds the serialization groups to use on this field / edge.
		Groups groups `json:"Groups,omitempty"`
		// MaxDepth tells the generator the maximum depth of this field when there is a cycle possible.
		MaxDepth uint
		// Validation holds the struct tags to use for github.com/go-playground/validator/v10. Used when no specific
		// validation tags are given in CreateValidation or UpdateValidation.
		Validation string
		// CreateValidation holds the struct tags to use for github.com/go-playground/validator/v10
		// when creating a new model.
		CreateValidation string
		// UpdateValidation holds the struct tags to use for github.com/go-playground/validator/v10
		// when updating an existing model.
		UpdateValidation string
	}
)

func CreateGroups(gs ...string) SchemaAnnotation {
	return SchemaAnnotation{CreateGroups: gs}
}

func ReadGroups(gs ...string) SchemaAnnotation {
	return SchemaAnnotation{ReadGroups: gs}
}

func UpdateGroups(gs ...string) SchemaAnnotation {
	return SchemaAnnotation{UpdateGroups: gs}
}

func ListGroups(gs ...string) SchemaAnnotation {
	return SchemaAnnotation{ListGroups: gs}
}

func Groups(gs ...string) Annotation {
	return Annotation{Groups: gs}
}

func MaxDepth(d uint) Annotation {
	return Annotation{MaxDepth: d}
}

func Validation(v string) Annotation {
	return Annotation{Validation: v}
}

func CreateValidation(v string) Annotation {
	return Annotation{CreateValidation: v}
}

func UpdateValidation(v string) Annotation {
	return Annotation{UpdateValidation: v}
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
	if ant.Validation != "" {
		a.Validation = ant.Validation
	}
	if ant.CreateValidation != "" {
		a.CreateValidation = ant.CreateValidation
	}
	if ant.UpdateValidation != "" {
		a.UpdateValidation = ant.UpdateValidation
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

// ValidationTags returns the tags to use for the given action.
func (a Annotation) ValidationTags(action string) string {
	if action == "create" && a.CreateValidation != "" {
		return a.CreateValidation
	}
	if action == "update" && a.UpdateValidation != "" {
		return a.UpdateValidation
	}
	return a.Validation
}

var (
	_ schema.Annotation = (*Annotation)(nil)
	_ schema.Annotation = (*SchemaAnnotation)(nil)
)
