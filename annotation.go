package elk

import "encoding/json"

type (
	// SchemaAnnotation annotates an entity with metadata for templates.
	SchemaAnnotation struct {
		// Skip tells the generator to skip this model.
		Skip bool `json:"Skip,omitempty"`
		// CreateGroups holds the serializations groups to use on the create handler.
		CreateGroups Groups `json:"CreateGroups,omitempty"`
		// ReadGroups holds the serializations groups to use on the read handler.
		ReadGroups Groups `json:"ReadGroups,omitempty"`
		// UpdateGroups holds the serializations groups to use on the update handler.
		UpdateGroups Groups `json:"UpdateGroups,omitempty"`
		// DeleteGroups holds the serializations groups to use on the delete handler.
		DeleteGroups Groups `json:"DeleteGroups,omitempty"`
	}
	// Annotation annotates fields and edges with metadata for templates.
	Annotation struct {
		// Skip tells the generator to skip this field / edge.
		Skip bool `json:"Skip,omitempty"`
		// Groups holds the serialization groups to use on this field / edge.
		Groups Groups `json:"Groups,omitempty"`
		// MaxDepth tells the generator the maximum depth of this field when there is a cycle possible.
		MaxDepth uint
	}
)

// Name implements ent.Annotation interface.
func (SchemaAnnotation) Name() string {
	return "ElkSchema"
}

// Name implements ent.Annotation interface.
func (Annotation) Name() string {
	return "Elk"
}

// Decode from ent.
func (a *Annotation) Decode(o interface{}) error {
	buf, err := json.Marshal(o)
	if err != nil {
		return err
	}

	return json.Unmarshal(buf, a)
}

// Decode from ent.
func (a *SchemaAnnotation) Decode(o interface{}) error {
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
