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
		vs, err := responseViews(g)
		if err != nil {
			return err
		}
		for _, v := range vs {
			n, err := v.ViewName()
			if err != nil {
				return err
			}
			ns = append(ns, n, n+"s")
		}
		for _, n := range g.Nodes {
			// Add the request structs used to deserialize request bodies.
			ns = append(ns,
				n.Name+"CreateRequest",
				n.Name+"UpdateRequest",
			)
		}

		// Run the easyjson generator.
		// TODO: Use ExtensionOptions to configure easyjson-options:
		//  - NoStdMarshalers
		//  - SnakeCase
		//  - LowerCamelCase
		//  - OmitEmpty
		//  - DisallowUnknownFields
		//  - SkipMemberNameUnescaping
		return (&bootstrap.Generator{
			PkgPath:         filepath.Join(g.Package, "http"),
			PkgName:         "http",
			Types:           ns,
			NoStdMarshalers: true,
			OutName:         filepath.Join(g.Config.Target, "http", "easyjson.go"),
		}).Run()
	})
}
