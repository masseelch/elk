package gen

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/masseelch/elk/internal"
)

func Generate(c *Config) error {
	cfg := &gen.Config{
		Target: c.Target,
		Package: c.Package,
		Templates: []*gen.Template{
			gen.MustParse(gen.NewTemplate("").Parse(string(internal.MustAsset("sheriff.tpl")))),
		},
	}

	if len(c.Templates) > 0 {
		return entc.Generate(c.Source, cfg, entc.TemplateFiles(c.Templates...))
	}

	return entc.Generate(c.Source, cfg)
}
