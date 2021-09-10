package elk

import (
	"encoding/json"
	"entgo.io/ent/schema"
	"github.com/masseelch/elk/serialization"
	"github.com/masseelch/elk/spec"
)

// SchemaAnnotation annotates an entity with metadata for templates.
type SchemaAnnotation struct {
	// CreatePolicy defines if a creation handler should be generated.
	CreatePolicy Policy
	// ReadPolicy defines if a read handler should be generated.
	ReadPolicy Policy
	// UpdatePolicy defines if an update handler should be generated.
	UpdatePolicy Policy
	// DeletePolicy defines if a delete handler should be generated.
	DeletePolicy Policy
	// ListPolicy defines if a list handler should be generated.
	ListPolicy Policy
	// CreateGroups holds the serializations groups to use on the creation handler.
	CreateGroups serialization.Groups
	// ReadGroups holds the serializations groups to use on the read handler.
	ReadGroups serialization.Groups
	// UpdateGroups holds the serializations groups to use on the update handler.
	UpdateGroups serialization.Groups
	// ListGroups holds the serializations groups to use on the list handler.
	ListGroups serialization.Groups
	// CreateSecurity sets the security property of the operation in the generated OpenAPI Spec.
	CreateSecurity spec.Security
	// ReadSecurity sets the security property of the operation in the generated OpenAPI Spec.
	ReadSecurity spec.Security
	// UpdateSecurity sets the security property of the operation in the generated OpenAPI Spec.
	UpdateSecurity spec.Security
	// DeleteSecurity sets the security property of the operation in the generated OpenAPI Spec.
	DeleteSecurity spec.Security
	// ListSecurity sets the security property of the operation in the generated OpenAPI Spec.
	ListSecurity spec.Security
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
	if ant.CreatePolicy != None {
		a.CreatePolicy = ant.CreatePolicy
	}
	if ant.ReadPolicy != None {
		a.ReadPolicy = ant.ReadPolicy
	}
	if ant.UpdatePolicy != None {
		a.UpdatePolicy = ant.UpdatePolicy
	}
	if ant.DeletePolicy != None {
		a.DeletePolicy = ant.DeletePolicy
	}
	if ant.ListPolicy != None {
		a.ListPolicy = ant.ListPolicy
	}
	if ant.CreateSecurity != nil {
		a.CreateSecurity = ant.CreateSecurity
	}
	if ant.ReadSecurity != nil {
		a.ReadSecurity = ant.ReadSecurity
	}
	if ant.UpdateSecurity != nil {
		a.UpdateSecurity = ant.UpdateSecurity
	}
	if ant.DeleteSecurity != nil {
		a.DeleteSecurity = ant.DeleteSecurity
	}
	if ant.ListSecurity != nil {
		a.ListSecurity = ant.ListSecurity
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

// CreateGroups returns a creation groups schema-annotation.
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

// CreatePolicy returns a creation policy schema-annotation.
func CreatePolicy(p Policy) SchemaAnnotation {
	return SchemaAnnotation{CreatePolicy: p}
}

// ReadPolicy returns a read policy schema-annotation.
func ReadPolicy(p Policy) SchemaAnnotation {
	return SchemaAnnotation{ReadPolicy: p}
}

// UpdatePolicy returns an update policy schema-annotation.
func UpdatePolicy(p Policy) SchemaAnnotation {
	return SchemaAnnotation{UpdatePolicy: p}
}

// DeletePolicy returns a delete policy schema-annotation.
func DeletePolicy(p Policy) SchemaAnnotation {
	return SchemaAnnotation{DeletePolicy: p}
}

// ListPolicy returns a list policy schema-annotation.
func ListPolicy(p Policy) SchemaAnnotation {
	return SchemaAnnotation{ListPolicy: p}
}

// SchemaPolicy returns a schema-annotation with all operation-policies set to the given one.
func SchemaPolicy(p Policy) SchemaAnnotation {
	return SchemaAnnotation{
		CreatePolicy: p,
		ReadPolicy:   p,
		UpdatePolicy: p,
		DeletePolicy: p,
		ListPolicy:   p,
	}
}

// SchemaSecurity sets the given security on all schema operations.
func SchemaSecurity(s spec.Security) SchemaAnnotation {
	return SchemaAnnotation{
		CreateSecurity: s,
		ReadSecurity:   s,
		UpdateSecurity: s,
		DeleteSecurity: s,
		ListSecurity:   s,
	}
}

// CreateSecurity returns a create-security schema-annotation.
func CreateSecurity(s spec.Security) SchemaAnnotation {
	return SchemaAnnotation{CreateSecurity: s}
}

// ReadSecurity returns a read-security schema-annotation.
func ReadSecurity(s spec.Security) SchemaAnnotation {
	return SchemaAnnotation{ReadSecurity: s}
}

// UpdateSecurity returns an update-security schema-annotation.
func UpdateSecurity(s spec.Security) SchemaAnnotation {
	return SchemaAnnotation{UpdateSecurity: s}
}

// DeleteSecurity returns a delete-security schema-annotation.
func DeleteSecurity(s spec.Security) SchemaAnnotation {
	return SchemaAnnotation{DeleteSecurity: s}
}

// ListSecurity returns a list-security schema-annotation.
func ListSecurity(s spec.Security) SchemaAnnotation {
	return SchemaAnnotation{ListSecurity: s}
}

// Annotation annotates fields and edges with metadata for templates.
type Annotation struct {
	// Groups holds the serialization groups to use on this field / edge.
	Groups serialization.Groups
	// MaxDepth tells the generator the maximum depth of this field when there is a cycle possible.
	MaxDepth uint
	// Expose defines if a read/list for this edge should be generated.
	Expose Policy
	// OpenAPI spec example value on schema fields.
	Example interface{}
	// OpenAPI security object for the read/list operation on this edge.
	Security spec.Security
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
	if ant.Expose != None {
		a.Expose = ant.Expose
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

// EnsureDefaults ensures defaults are set.
func (a *Annotation) EnsureDefaults() {
	if a.MaxDepth == 0 {
		a.MaxDepth = 1
	}
}

// Groups returns a groups annotation.
func Groups(gs ...string) Annotation {
	return Annotation{Groups: gs}
}

// MaxDepth returns a max depth annotation.
func MaxDepth(d uint) Annotation {
	return Annotation{MaxDepth: d}
}

// ExposeEdge returns a Expose annotation.
func ExposeEdge() Annotation {
	return Annotation{Expose: Expose}
}

// ExcludeEdge returns a Exclude annotation.
func ExcludeEdge() Annotation {
	return Annotation{Expose: Exclude}
}

// Example returns an example annotation.
func Example(v interface{}) Annotation {
	return Annotation{Example: v}
}

var (
	_ schema.Annotation = (*Annotation)(nil)
	_ schema.Annotation = (*SchemaAnnotation)(nil)
)
