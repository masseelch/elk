package gen

import (
	"github.com/facebook/ent/entc"
	"github.com/facebook/ent/entc/gen"
	"github.com/masseelch/elk/internal"
)

func Generate(source string, target string) error {
	return entc.Generate(source, &gen.Config{
		Target: target,
		Templates: []*gen.Template{
			gen.MustParse(gen.NewTemplate("").Parse(string(internal.MustAsset("sheriff.tpl")))),
		},
	})
}
