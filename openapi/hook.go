package openapi

import "entgo.io/ent/entc/gen"

func Hook(spec *Spec) gen.Hook {
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			return nil
		})
	}
}
