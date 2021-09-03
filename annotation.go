package elk

import (
	"encoding/json"
	"entgo.io/ent/schema"
	"github.com/masseelch/elk/policy"
	"github.com/masseelch/elk/serialization"
)

// SchemaAnnotation annotates an entity with metadata for templates.
type SchemaAnnotation struct {
	// ExposeCreate defines if a creation handler should be generated.
	CreatePolicy policy.Policy
	// ExposeRead defines if a read handler should be generated.
	ReadPolicy policy.Policy
	// ExposeUpdate defines if an update handler should be generated.
	UpdatePolicy policy.Policy
	// ExposeDelete defines if a delete handler should be generated.
	DeletePolicy policy.Policy
	// ExposeList defines if a list handler should be generated.
	ListPolicy policy.Policy
	// CreateGroups holds the serializations groups to use on the creation handler.
	CreateGroups serialization.Groups
	// ReadGroups holds the serializations groups to use on the read handler.
	ReadGroups serialization.Groups
	// UpdateGroups holds the serializations groups to use on the update handler.
	UpdateGroups serialization.Groups
	// ListGroups holds the serializations groups to use on the list handler.
	ListGroups serialization.Groups
	// // CreateSummary overrides default OpenAPI spec operation summary on create operation.
	// CreateSummary string
	// // CreateDescription overrides default OpenAPI spec operation description on create operation.
	// CreateDescription string
	// // ReadSummary overrides default OpenAPI spec operation summary on read operation.
	// ReadSummary string
	// // ReadDescription overrides default OpenAPI spec operation description on read operation.
	// ReadDescription string
	// // UpdateSummary overrides default OpenAPI spec operation summary on update operation.
	// UpdateSummary string
	// // UpdateDescription overrides default OpenAPI spec operation description on update operation.
	// UpdateDescription string
	// // DeleteSummary overrides default OpenAPI spec operation summary on delete operation.
	// DeleteSummary string
	// // DeleteDescription overrides default OpenAPI spec operation description on delete operation.
	// DeleteDescription string
	// // ListSummary overrides default OpenAPI spec operation summary on list operation.
	// ListSummary string
	// // ListDescription overrides default OpenAPI spec operation description on list operation.
	// ListDescription string
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
	if ant.CreatePolicy != policy.None {
		a.CreatePolicy = ant.CreatePolicy
	}
	if ant.ReadPolicy != policy.None {
		a.ReadPolicy = ant.ReadPolicy
	}
	if ant.UpdatePolicy != policy.None {
		a.UpdatePolicy = ant.UpdatePolicy
	}
	if ant.DeletePolicy != policy.None {
		a.DeletePolicy = ant.DeletePolicy
	}
	if ant.ListPolicy != policy.None {
		a.ListPolicy = ant.ListPolicy
	}
	// if ant.CreateSummary != "" {
	// 	a.CreateSummary = ant.CreateSummary
	// }
	// if ant.CreateDescription != "" {
	// 	a.CreateDescription = ant.CreateDescription
	// }
	// if ant.ReadSummary != "" {
	// 	a.ReadSummary = ant.ReadSummary
	// }
	// if ant.ReadDescription != "" {
	// 	a.ReadDescription = ant.ReadDescription
	// }
	// if ant.UpdateSummary != "" {
	// 	a.UpdateSummary = ant.UpdateSummary
	// }
	// if ant.UpdateDescription != "" {
	// 	a.UpdateDescription = ant.UpdateDescription
	// }
	// if ant.DeleteSummary != "" {
	// 	a.DeleteSummary = ant.DeleteSummary
	// }
	// if ant.DeleteDescription != "" {
	// 	a.DeleteDescription = ant.DeleteDescription
	// }
	// if ant.ListSummary != "" {
	// 	a.ListSummary = ant.ListSummary
	// }
	// if ant.ListDescription != "" {
	// 	a.ListDescription = ant.ListDescription
	// }
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

type PolicyConfig uint

const (
	Create PolicyConfig = iota
	Read
	Update
	Delete
	List
)

// Expose enables all CRUD-operations on a schema.
func Expose(c ...PolicyConfig) SchemaAnnotation {
	// If no config is given set all to policy.Expose.
	if len(c) == 0 {
		return SchemaAnnotation{
			CreatePolicy: policy.Expose,
			ReadPolicy:   policy.Expose,
			UpdatePolicy: policy.Expose,
			DeletePolicy: policy.Expose,
			ListPolicy:   policy.Expose,
		}
	}
	// If a config is given only set those to policy.Expose that are requested.
	s := SchemaAnnotation{}
	for _, c := range c {
		switch c {
		case Create:
			s.CreatePolicy = policy.Expose
		case Read:
			s.ReadPolicy = policy.Expose
		case Update:
			s.UpdatePolicy = policy.Expose
		case Delete:
			s.DeletePolicy = policy.Expose
		case List:
			s.ListPolicy = policy.Expose
		}
	}
	return s
}

// Exclude disables all CRUD-operations on a schema.
func Exclude(c ...PolicyConfig) SchemaAnnotation {
	// If no config is given set all to policy.Expose.
	if len(c) == 0 {
		return SchemaAnnotation{
			CreatePolicy: policy.Exclude,
			ReadPolicy:   policy.Exclude,
			UpdatePolicy: policy.Exclude,
			DeletePolicy: policy.Exclude,
			ListPolicy:   policy.Exclude,
		}
	}
	// If a config is given only set those to policy.Expose that are requested.
	s := SchemaAnnotation{}
	for _, c := range c {
		switch c {
		case Create:
			s.CreatePolicy = policy.Exclude
		case Read:
			s.ReadPolicy = policy.Exclude
		case Update:
			s.UpdatePolicy = policy.Exclude
		case Delete:
			s.DeletePolicy = policy.Exclude
		case List:
			s.ListPolicy = policy.Exclude
		}
	}
	return s
}

// // CreateSummaryDescription adds OpenAPI spec operation summary and description to a create operation.
// func CreateSummaryDescription(sum, desc string) SchemaAnnotation {
// 	return SchemaAnnotation{
// 		CreateSummary:     sum,
// 		CreateDescription: desc,
// 	}
// }
//
// // ReadSummaryDescription adds OpenAPI spec operation summary and description to a read operation.
// func ReadSummaryDescription(sum, desc string) SchemaAnnotation {
// 	return SchemaAnnotation{
// 		ReadSummary:     sum,
// 		ReadDescription: desc,
// 	}
// }
//
// // UpdateSummaryDescription adds OpenAPI spec operation summary and description to an update operation.
// func UpdateSummaryDescription(sum, desc string) SchemaAnnotation {
// 	return SchemaAnnotation{
// 		UpdateSummary:     sum,
// 		UpdateDescription: desc,
// 	}
// }
//
// // DeleteSummaryDescription adds OpenAPI spec operation summary and description to a delete operation.
// func DeleteSummaryDescription(sum, desc string) SchemaAnnotation {
// 	return SchemaAnnotation{
// 		DeleteSummary:     sum,
// 		DeleteDescription: desc,
// 	}
// }
//
// // ListSummaryDescription adds OpenAPI spec operation summary and description to a list operation.
// func ListSummaryDescription(sum, desc string) SchemaAnnotation {
// 	return SchemaAnnotation{
// 		ListSummary:     sum,
// 		ListDescription: desc,
// 	}
// }

// Annotation annotates fields and edges with metadata for templates.
type Annotation struct {
	// Groups holds the serialization groups to use on this field / edge.
	Groups serialization.Groups
	// MaxDepth tells the generator the maximum depth of this field when there is a cycle possible.
	MaxDepth uint
	// Expose defines if a read/list for this edge should be generated.
	Expose policy.Policy
	// OpenAPI spec example value on schema fields.
	Example interface{}
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
	if ant.Expose != policy.None {
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

// ExposeEdge returns a policy.Expose annotation.
func ExposeEdge() Annotation {
	return Annotation{Expose: policy.Expose}
}

// ExcludeEdge returns a policy.Exclude annotation.
func ExcludeEdge() Annotation {
	return Annotation{Expose: policy.Exclude}
}

// Example returns an example annotation.
func Example(v interface{}) Annotation {
	return Annotation{Example: v}
}

var (
	_ schema.Annotation = (*Annotation)(nil)
	_ schema.Annotation = (*SchemaAnnotation)(nil)
)
