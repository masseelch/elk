package elk

import (
	"embed"
	"entgo.io/ent/entc/gen"
	"fmt"
	"github.com/stoewer/go-strcase"
	"text/template"
)

var (
	//go:embed template
	templateDir embed.FS
	// Funcs contains the extra template functions used by elk.
	Funcs = template.FuncMap{
		"contains":        contains,
		"edges":           edges,
		"filterEdges":     filterEdges,
		"filterNodes":     filterNodes,
		"kebab":           strcase.KebabCase,
		"needsValidation": needsValidation,
		"nodeOperations":  nodeOperations,
		"view":            newView,
		"views":           newViews,
		"stringSlice":     stringSlice,
		"xextend":         xextend,
		"zapField":        zapField,
	}
	// HTTPTemplate holds all templates for generating http handlers.
	HTTPTemplate = gen.MustParse(gen.NewTemplate("elk").Funcs(Funcs).ParseFS(templateDir, "template/http/*.tmpl"))
)

// filterNodes returns the nodes a handler for the given operation should be generated for.
func filterNodes(g *gen.Graph, op string) ([]*gen.Type, error) {
	c, err := config(g.Config)
	if err != nil {
		return nil, err
	}
	var filteredNodes []*gen.Type
	for _, n := range g.Nodes {
		var p Policy
		ant := &SchemaAnnotation{}
		// If no policies are given follow the global policy.
		if n.Annotations == nil || n.Annotations[ant.Name()] == nil {
			p = c.HandlerPolicy
		} else {
			if err := ant.Decode(n.Annotations[ant.Name()]); err != nil {
				return nil, err
			}
			switch op {
			case createOperation:
				p = ant.CreatePolicy
			case readOperation:
				p = ant.ReadPolicy
			case updateOperation:
				p = ant.UpdatePolicy
			case deleteOperation:
				p = ant.DeletePolicy
			case listOperation:
				p = ant.ListPolicy
			}
			// If the policy is policy.None follow the globally defined policy.
			if p == None {
				p = c.HandlerPolicy
			}
		}
		if p == Expose {
			filteredNodes = append(filteredNodes, n)
		}
	}
	return filteredNodes, nil
}

// filterEdges returns the edges a read/list handler should be generated for.
func filterEdges(n *gen.Type) ([]*gen.Edge, error) {
	c, err := config(n.Config)
	if err != nil {
		return nil, err
	}
	var filteredEdges []*gen.Edge
	for _, e := range n.Edges {
		var p Policy
		ant := &Annotation{}
		// If no policies are given follow the global policy.
		if e.Annotations == nil || e.Annotations[ant.Name()] == nil {
			p = c.HandlerPolicy
		} else {
			if err := ant.Decode(e.Annotations[ant.Name()]); err != nil {
				return nil, err
			}
			p = ant.Expose
			// If the policy is policy.None follow the globally defined policy.
			if p == None {
				p = c.HandlerPolicy
			}
		}
		if p == Expose {
			filteredEdges = append(filteredEdges, e)
		}
	}
	return filteredEdges, nil
}

// nodeOperations returns the list of operations to expose for this node.
func nodeOperations(n *gen.Type) ([]string, error) {
	c, err := config(n.Config)
	if err != nil {
		return nil, err
	}
	ops := []string{createOperation, readOperation, updateOperation, deleteOperation, listOperation}
	ant := &SchemaAnnotation{}
	// If no policies are given follow the global policy.
	if n.Annotations == nil || n.Annotations[ant.Name()] == nil {
		if c.HandlerPolicy == Expose {
			return ops, nil
		}
		return nil, nil
	} else {
		if err := ant.Decode(n.Annotations[ant.Name()]); err != nil {
			return nil, err
		}
		var ops []string
		if ant.CreatePolicy == Expose || (ant.CreatePolicy == None && c.HandlerPolicy == Expose) {
			ops = append(ops, createOperation)
		}
		if ant.ReadPolicy == Expose || (ant.ReadPolicy == None && c.HandlerPolicy == Expose) {
			ops = append(ops, readOperation)
		}
		if ant.UpdatePolicy == Expose || (ant.UpdatePolicy == None && c.HandlerPolicy == Expose) {
			ops = append(ops, updateOperation)
		}
		if ant.DeletePolicy == Expose || (ant.DeletePolicy == None && c.HandlerPolicy == Expose) {
			ops = append(ops, deleteOperation)
		}
		if ant.ListPolicy == Expose || (ant.ListPolicy == None && c.HandlerPolicy == Expose) {
			ops = append(ops, listOperation)
		}
		return ops, nil
	}
}

// needsValidation returns if a type needs validation because there is some defined on one of its fields.
func needsValidation(n *gen.Type) bool {
	for _, f := range n.Fields {
		if f.Validators > 0 {
			return true
		}
	}
	return false
}

// contains checks if a string slice contains the given value.
func contains(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}

func stringSlice(src []interface{}) []string {
	if len(src) == 0 {
		return nil
	}
	dst := make([]string, len(src))
	for i, v := range src {
		dst[i] = v.(string)
	}
	return dst
}

// edgeScope wraps the Edge object with extended scope.
type edgeScope struct {
	Edge
	Scope map[interface{}]interface{}
}

// xextend extends the parent block with a KV pairs. Stolen from entgo.io/ent/entc/gen/func.go.
//
//	{{ with $scope := xextend $ "key" "value" }}
//		{{ template "setters" $scope }}
//	{{ end}}
//
func xextend(v interface{}, kv ...interface{}) (interface{}, error) {
	scope := make(map[interface{}]interface{})
	if len(kv)%2 != 0 {
		return nil, fmt.Errorf("invalid number of parameters: %d", len(kv))
	}
	for i := 0; i < len(kv); i += 2 {
		scope[kv[i]] = kv[i+1]
	}
	switch v := v.(type) {
	case Edge:
		return &edgeScope{Edge: v, Scope: scope}, nil
	case *edgeScope:
		for k := range v.Scope {
			scope[k] = v.Scope[k]
		}
		return &edgeScope{Edge: v.Edge, Scope: scope}, nil
	default:
		return nil, fmt.Errorf("invalid type for xextend: %T", v)
	}
}

// zapField returns the method-name to use for logging fields.
func zapField(f *gen.Field) (string, error) {
	switch {
	case f.IsString():
		return "String", nil
	case f.IsInt():
		return "Int", nil
	}
	return "", fmt.Errorf("elk: invalid ID-Type %q", f.Type.String())
}
