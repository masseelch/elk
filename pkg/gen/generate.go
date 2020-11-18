package gen

import (
	"github.com/facebook/ent/entc"
	"github.com/facebook/ent/entc/gen"
	"github.com/masseelch/elk/internal"
)

func Generate(c *Config) error {
	return entc.Generate(c.Source, &gen.Config{
		Target: c.Target,
		Templates: []*gen.Template{
			gen.MustParse(gen.NewTemplate("").Parse(string(internal.MustAsset("sheriff.tpl")))),
		},
	})
}
