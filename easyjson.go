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
			ns = append(ns,
				// Add all four operation with response data for every node.
				n.Name+"CreateResponse",
				n.Name+"ReadResponse",
				n.Name+"UpdateResponse",
				n.Name+"ListResponse",
				// Add the request structs used to deserialize request bodies.
				n.Name+"CreateRequest",
				n.Name+"UpdateRequest",
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
