package gen

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/masseelch/elk/internal"
)

func Generate(c *Config) error {
	return entc.Generate(c.Source, &gen.Config{
		Target: c.Target,
		Package: c.Package,
		Templates: []*gen.Template{
			gen.MustParse(gen.NewTemplate("").Parse(string(internal.MustAsset("sheriff.tpl")))),
		},
	})
}
