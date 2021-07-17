package elk

import (
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
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
	}
	// TemplateFuncs contains the extra template functions used by elk.
	TemplateFuncs = template.FuncMap{
		"edgesToLoad":    edgesToLoad,
		"dartType":       dartType,
		"kebab":          strcase.KebabCase,
		"validationTags": validationTags,
	}
	// TypeMappings contains the string representation in dart for a given go type.
	TypeMappings = map[string]string{
		"invalid":   "dynamic",
		"bool":      "bool",
		"time.Time": "DateTime",
		// "JSON":    "Map<String, dynamic>",
		// "UUID":    "String",
		// "bytes":   "dynamic",
		"enum":     "String",
		"string":   "String",
		"int":      "int",
		"int8":     "int",
		"int16":    "int",
		"int32":    "int",
		"int64":    "int",
		"uint":     "int",
		"uint8":    "int",
		"uint16":   "int",
		"uint32":   "int",
		"uint64":   "int",
		"float32":  "double",
		"float64":  "double",
		"[]string": "List<String>",
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

// dartType returns the dart type for a given schema field.
func dartType(f *field.TypeInfo) string {
	// If there is an entry in the map use it.
	if t, ok := TypeMappings[f.String()]; ok {
		return t
	}
	// Try to guess the type. Returns dynamic in an invalid case.
	return TypeMappings[f.Type.String()]
}
