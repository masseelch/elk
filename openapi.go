package elk

import (
	"encoding/json"
	"entgo.io/ent/entc/gen"
	"fmt"
	oas "github.com/masseelch/elk/spec"
)

var (
	// _base64    = &spec.Type{"string", "byte"}
	// _uint8List = &spec.Type{"string", "binary"}
	// _date      = &spec.Type{"string", "date"}
	// _sensitive = &spec.Type{"string", "password"}
	_int32    = &oas.Type{Type: "integer", Format: "int32"}
	_int64    = &oas.Type{Type: "integer", Format: "int64"}
	_float    = &oas.Type{Type: "number", Format: "float"}
	_double   = &oas.Type{Type: "number", Format: "double"}
	_string   = &oas.Type{Type: "string"}
	_bool     = &oas.Type{Type: "boolean"}
	_dateTime = &oas.Type{Type: "string", Format: "date-time"}
	oasTypes  = map[string]*oas.Type{
		"bool":      _bool,
		"time.Time": _dateTime,
		"enum":      _string,
		"string":    _string,
		"int":       _int32,
		"int8":      _int32,
		"int16":     _int32,
		"int32":     _int32,
		"uint":      _int32,
		"uint8":     _int32,
		"uint16":    _int32,
		"uint32":    _int32,
		"int64":     _int64,
		"uint64":    _int64,
		"float32":   _float,
		"float64":   _double,
	}
)

// SpecGenerator TODO
func SpecGenerator(spec *oas.Spec) gen.Hook {
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			// Let ent create all the files.
			if err := next.Generate(g); err != nil {
				return err
			}
			// Ensure spec is ready to receive data.
			spec.WarmUp()
			// Loop over every node and add its routes to the spec (including all views).
			for _, n := range g.Nodes {
				// The schema fields.
				fields := make(oas.Fields, len(n.Fields))
				for _, f := range n.Fields {
					t, ok := oasTypes[f.Type.String()]
					if !ok {
						return fmt.Errorf("no OAS-type exists for %q", f.Type.String())
					}
					var e interface{}
					a := Annotation{}
					if f.Annotations != nil && f.Annotations[a.Name()] != nil {
						if err := a.Decode(f.Annotations[a.Name()]); err != nil {
							return err
						}
						e = a.Example
					}
					fields[f.Name] = oas.Field{
						Required: !f.Optional,
						Type:     *t,
						Example:  e,
					}
				}
				// The schema edges.
				// edges := make(oas.Edges, len(n.Edges))
				// for _, e := range n.Edges {
				//
				// }
				spec.Components.Schemas[n.Name] = oas.Schema{
					Fields: fields,
				}
			}
			return dump(spec)
		})
	}
}

func dump(spec *oas.Spec) error {
	b, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
