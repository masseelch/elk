package elk

import (
	"entgo.io/ent/entc/gen"
	"fmt"
	"github.com/masseelch/elk/internal"
	"github.com/stoewer/go-strcase"
	"text/template"
)

//go:generate go run github.com/go-bindata/go-bindata/go-bindata -o=internal/bindata.go -pkg=internal -modtime=1 ./template/...

const (
	actionCreate = "create"
	actionRead   = "read"
	actionUpdate = "update"
	actionList   = "list"
)

var (
	// HTTPTemplates holds all templates for generating http handlers.
	HTTPTemplates = []*gen.Template{
		parse("template/http/handler.tmpl"),
		parse("template/http/create.tmpl"),
		parse("template/http/read.tmpl"),
		parse("template/http/update.tmpl"),
		parse("template/http/delete.tmpl"),
		parse("template/http/list.tmpl"),
		parse("template/http/relations.tmpl"),
		parse("template/http/request.tmpl"),
		parse("template/http/response.tmpl"),
		parse("template/http/helpers.tmpl"),
		parse("template/http/import.tmpl"),
	}
	// TemplateFuncs contains the extra template functions used by elk.
	TemplateFuncs = template.FuncMap{
		"edgesToLoad":             edgesToLoad,
		"fieldValidationRequired": fieldValidationRequired,
		"kebab":                   strcase.KebabCase,
		"needsSerialization":      needsSerialization,
		"stringSlice":             stringSlice,
		"validationTags":          validationTags,
		"xextend":                 xextend,
	}
)

func parse(path string) *gen.Template {
	return gen.MustParse(gen.NewTemplate(path).
		Funcs(TemplateFuncs).
		Parse(string(internal.MustAsset(path))))
}

// validationTags extracts the validation tags to use for the given action / method.
func validationTags(a interface{}, m string) string {
	if a == nil {
		return ""
	}
	an := Annotation{}
	if err := an.Decode(a); err != nil {
		return ""
	}
	if m == "create" && an.CreateValidation != "" {
		return an.CreateValidation
	}
	if m == "update" && an.UpdateValidation != "" {
		return an.UpdateValidation
	}
	return an.Validation
}

// fieldValidationRequired returns if a type needs validation because there is some defined on one of its fields.
func fieldValidationRequired(n *gen.Type) bool {
	for _, f := range n.Fields {
		if f.Validators > 0 {
			return true
		}
	}

	return false
}

// stringSlice casts a given []interface{} to []string.
func stringSlice(is interface{}) []string {
	switch is := is.(type) {
	case []interface{}:
		ss := make([]string, len(is))
		for i, v := range is {
			ss[i] = v.(string)
		}
		return ss
	case []string:
		return is
	default:
		return nil
	}
}

// needsSerialization checks if a given field  is to be serialized according to its annotations and the requested
// groups.
func needsSerialization(a interface{}, g groups) (bool, error) {
	// If no groups are requested or the field has no groups defined render the field.
	if a == nil || len(g) == 0 {
		return true, nil
	}
	// If there are groups given check if the groups match the requested ones.
	an := Annotation{}
	if err := an.Decode(a); err != nil {
		return false, err
	}
	// If no groups are given on the field default is to include it in the output.
	if len(an.Groups) == 0 {
		return true, nil
	}
	return g.Match(an.Groups), nil
}

// graphScope wraps the Graph object with extended scope.
type edgeToLoadScope struct {
	EdgeToLoad
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
	case EdgeToLoad:
		return &edgeToLoadScope{EdgeToLoad: v, Scope: scope}, nil
	case *edgeToLoadScope:
		for k := range v.Scope {
			scope[k] = v.Scope[k]
		}
		return &edgeToLoadScope{EdgeToLoad: v.EdgeToLoad, Scope: scope}, nil
	default:
		return nil, fmt.Errorf("invalid type for xextend: %T", v)
	}
}
