package openapi

import (
	"entgo.io/ent/entc/gen"
	"github.com/masseelch/elk/openapi/spec"
)

// Generator TODO
func Generator(s *spec.Spec) gen.Hook {
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			// Let ent create all the files.
			if err := next.Generate(g); err != nil {
				return err
			}
			// Ensure spec is ready to receive data.
			s.WarmUp()
			// Loop over every node and add it with its routes to the spec.
			for _, n := range g.Nodes {
				fields := make(map[string]spec.Field, len(n.Fields))
				for _, f := range n.Fields {
					fields[f.Name] = spec.Field{
						Required: !f.Optional,
						Type:     "",  // TODO: map of Go-Primitives to OAS-Types
						Format:   "",  // TODO: same as above
						Example:  nil, // TODO: Read from annotation or empty
					}
				}
				s.Components.Schemas[n.Name] = spec.Schema{
					Fields: fields,
				}
			}
			return nil
		})
	}
}
