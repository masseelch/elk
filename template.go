package elk

import (
	"entgo.io/ent/entc/gen"
	"github.com/masseelch/elk/internal"
	"github.com/stoewer/go-strcase"
	"text/template"
)

//go:generate go run github.com/go-bindata/go-bindata/go-bindata -o=internal/bindata.go -pkg=internal -modtime=1 ./template/...

const (
	actionCreate = "create"
	actionRead   = "read"
	actionUpdate = "update"
	actionDelete = "delete"
	actionList   = "list"
)

var (
	// HTTPTemplates holds all templates for generating http handlers.
	HTTPTemplates = []*gen.Template{
		parse("template/http/handler.tmpl"),
		parse("template/http/helpers.tmpl"),
		parse("template/http/create.tmpl"),
		parse("template/http/read.tmpl"),
		parse("template/http/update.tmpl"),
		parse("template/http/delete.tmpl"),
		parse("template/http/list.tmpl"),
		parse("template/http/relations.tmpl"),
	}
	// TemplateFuncs contains the extra template functions used by elk.
	TemplateFuncs = template.FuncMap{
		"edgesToLoad":             edgesToLoad,
		"kebab":                   strcase.KebabCase,
		"stringSlice":             stringSlice,
		"validationTags":          validationTags,
		"fieldValidationRequired": fieldValidationRequired,
	}
)

func parse(path string) *gen.Template {
	return gen.MustParse(gen.NewTemplate(path).
		Funcs(gen.Funcs).
		Funcs(TemplateFuncs).
		Parse(string(internal.MustAsset(path))))
}

func validationTags(a gen.Annotations, m string) string {
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

func fieldValidationRequired(n *gen.Type) bool {
	for _, f := range n.Fields {
		if f.Validators > 0 {
			return true
		}
	}

	return false
}

// stringSlice casts a given []interface{} to []string.
func stringSlice(is []interface{}) []string {
	ss := make([]string, len(is))
	for i, v := range is {
		ss[i] = v.(string)
	}
	return ss
}
