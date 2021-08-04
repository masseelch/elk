package elk

import (
	"entgo.io/ent/entc/gen"
	"fmt"
	"github.com/masseelch/elk/internal"
	"github.com/stoewer/go-strcase"
	"regexp"
	"strings"
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
		parse("template/http/create.tmpl"),
		parse("template/http/delete.tmpl"),
		parse("template/http/handler.tmpl"),
		parse("template/http/helpers.tmpl"),
		parse("template/http/list.tmpl"),
		parse("template/http/read.tmpl"),
		parse("template/http/request.tmpl"),
		parse("template/http/relations.tmpl"),
		parse("template/http/response.tmpl"),
		parse("template/http/update.tmpl"),
	}
	// TemplateFuncs contains the extra template functions used by elk.
	TemplateFuncs = template.FuncMap{
		"edgesToLoad":        edgesToLoad,
		"kebab":              strcase.KebabCase,
		"needsSerialization": needsSerialization,
		"needsValidation":    needsValidation,
		"renderTypes":        renderTypes,
		"stringSlice":        stringSlice,
		"validationTags":     validationTags,
		"viewName":           viewName,
		"xappend":            xappend,
		"xextend":            xextend,
	}
	// SpecialCharsRegex to match all special chars.
	SpecialCharsRegex = regexp.MustCompile("[^0-9a-zA-Z_]+")
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

// needsValidation returns if a type needs validation for a given request type.
func needsValidation(n *gen.Type, m string) bool {
	an := Annotation{}.Name()
	for _, f := range n.Fields {
		if validationTags(f.Annotations[an], m) != "" {
			return true
		}
	}
	for _, e := range n.Edges {
		if validationTags(e.Annotations[an], m) != "" {
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

type renderType struct {
	*gen.Type
	Groups groups
}

// renderTypes creates a map of
func renderTypes(g *gen.Graph) (map[string]renderType, error) {
	m := make(map[string]renderType, 0)
	for _, n := range g.Nodes {
		a := SchemaAnnotation{}
		if err := a.Decode(n.Annotations[a.Name()]); err != nil {
			return nil, err
		}
		for _, gs := range [][]string{a.CreateGroups, a.ReadGroups, a.UpdateGroups, a.ListGroups} {
			m[viewName(n, gs)] = renderType{
				Type:   n,
				Groups: gs,
			}
		}
	}

	return m, nil
}

// viewName composes a name for the struct to use if the given groups a requested on the given schema.
func viewName(n *gen.Type, gs groups) string {
	if len(gs) == 0 {
		return n.Name + "Response"
	} else {
		return n.Name + "_" + SpecialCharsRegex.ReplaceAllString(strings.Join(gs, "_"), "") + "_Response"
	}
}

func xappend(xs []interface{}, ys ...interface{}) []interface{} {
	return append(xs, ys...)
}

// graphScope wraps the Graph object with extended scope.
type edgeToLoadScope struct {
	EdgeToLoad
	Scope map[interface{}]interface{}
}

// xextend extends the parent block with a KV pairs. Stolen from entgo.io/ent/entc/gen/func.go.
//
//	{{ with $scope := extend $ "key" "value" }}
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
