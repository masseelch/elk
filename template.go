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
		"edges":           edges,
		"kebab":           strcase.KebabCase,
		"needsValidation": needsValidation,
		"view":            view,
		"views":           views,
		"stringSlice":     stringSlice,
		"xextend":         xextend,
	}
	// HTTPTemplate holds all templates for generating http handlers.
	HTTPTemplate = gen.MustParse(gen.NewTemplate("elk").Funcs(Funcs).ParseFS(templateDir, "template/http/*.tmpl"))
)

// needsValidation returns if a type needs validation because there is some defined on one of its fields.
func needsValidation(n *gen.Type) bool {
	for _, f := range n.Fields {
		if f.Validators > 0 {
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
