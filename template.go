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
)

var (
	// HTTPTemplates holds all templates for generating http handlers.
	HTTPTemplates = []*gen.Template{
		parse("template/http/handler.tmpl"),
		parse("template/http/create.tmpl"),
		parse("template/http/read.tmpl"),
	}
	// TemplateFuncs contains the extra template functions used by elk.
	TemplateFuncs = template.FuncMap{
		"edgesToLoad": edgesToLoad,
		"kebab":       strcase.KebabCase,
	}
)

func parse(path string) *gen.Template {
	return gen.MustParse(gen.NewTemplate(path).
		Funcs(gen.Funcs).
		Funcs(TemplateFuncs).
		Parse(string(internal.MustAsset(path))))
}
