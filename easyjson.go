package elk

import (
	"entgo.io/ent/entc/gen"
	"github.com/mailru/easyjson/bootstrap"
	"path/filepath"
)

func GenerateEasyJSON(next gen.Generator) gen.Generator {
	return gen.GenerateFunc(func(g *gen.Graph) error {
		// Let ent create all the files.
		if err := next.Generate(g); err != nil {
			return err
		}

		// We want to render every response struct created with easyjson.
		var ns []string
		for _, n := range g.Nodes {
			// Add all four operation with response data for every node.
			ns = append(ns,
				n.Name+"CreateResponse",
				n.Name+"ReadResponse",
				n.Name+"UpdateResponse",
				n.Name+"ListResponse",
			)
		}

		// Run the easyjson generator.
		return (&bootstrap.Generator{
			PkgPath:               filepath.Join(g.Package, "http"),
			PkgName:               "http",
			Types:                 ns,
			NoStdMarshalers:       true,
			DisallowUnknownFields: true,
			OutName:               filepath.Join(g.Config.Target, "http", "easyjson.go"),
		}).Run()
	})
}
