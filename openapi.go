package elk

import (
	"encoding/json"
	"entgo.io/ent/entc/gen"
	"fmt"
	"github.com/go-openapi/inflect"
	"github.com/masseelch/elk/spec"
	"github.com/stoewer/go-strcase"
	"io"
	"net/http"
	"os"
	"strconv"
)

var (
	// _base64    = &spec.Type{"string", "byte"}
	// _uint8List = &spec.Type{"string", "binary"}
	// _date      = &spec.Type{"string", "date"}
	// _sensitive = &spec.Type{"string", "password"}
	_int32    = &spec.Type{Type: "integer", Format: "int32"}
	_int64    = &spec.Type{Type: "integer", Format: "int64"}
	_float    = &spec.Type{Type: "number", Format: "float"}
	_double   = &spec.Type{Type: "number", Format: "double"}
	_string   = &spec.Type{Type: "string"}
	_bool     = &spec.Type{Type: "boolean"}
	_dateTime = &spec.Type{Type: "string", Format: "date-time"}
	oasTypes  = map[string]*spec.Type{
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
	rules = inflect.NewDefaultRuleset()
)

type (
	// Generator is the interface that wraps the Generate method.
	Generator interface {
		// Generate edits the given OpenAPI spec.
		Generate(*spec.Spec) error
	}
	// The GenerateFunc type is an adapter to allow the use of ordinary
	// function as Generator. If f is a function with the appropriate signature,
	// GenerateFunc(f) is a Generator that calls f.
	GenerateFunc func(*spec.Spec) error
	// Hook defines the "spec generate middleware".
	Hook func(Generator) Generator
)

// Generate calls f(s).
func (f GenerateFunc) Generate(s *spec.Spec) error {
	return f(s)
}

// SpecGenerator TODO
func (e *Extension) SpecGenerator(out io.Writer) gen.Hook {
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			// Let ent create all the files.
			// if err := next.Generate(g); err != nil {
			// 	return err
			// }
			// Start the Generator chain.
			var chain Generator = generate(g)
			// Add user hooks to chain.
			for i := len(e.specHooks) - 1; i >= 0; i-- {
				chain = e.specHooks[i](chain)
			}
			// Create a fresh spec.
			s := initSpec()
			// Run the generators.
			if err := chain.Generate(s); err != nil {
				return err
			}
			// Dump the spec.
			return dump(out, s)
		})
	}
}

// initSpec returns an empty spec ready to receive data.
func initSpec() *spec.Spec {
	return &spec.Spec{
		Info: &spec.Info{
			Title:       "Ent Schema API",
			Description: "This is an auto generated API description made out of an Ent schema definition",
			Version:     "0.0.0",
		},
		Components: spec.Components{
			Schemas:    make(map[string]*spec.Schema),
			Responses:  make(map[string]spec.Response),
			Parameters: make(map[string]spec.Parameter),
		},
	}
}

// generate is the default Generator to fill a given spec.
func generate(g *gen.Graph) GenerateFunc {
	return func(s *spec.Spec) error {
		// Add all views to the schemas.
		if err := viewSchemas(g, s); err != nil {
			return err
		}
		// Create the paths.
		if err := paths(g, s); err != nil {
			return err
		}
		return nil
	}
}

// viewSchemas adds all views to the specs schemas.
func viewSchemas(g *gen.Graph, s *spec.Spec) error {
	vs, err := newViews(g)
	if err != nil {
		return err
	}
	// Create a schema for every view.
	for n, v := range vs {
		fs := make(spec.Fields, len(v.Fields))
		// We can already add the schema fields.
		for _, f := range v.Fields {
			sf, err := newField(f)
			if err != nil {
				return err
			}
			fs[f.Name] = sf
		}
		s.Components.Schemas[n] = &spec.Schema{
			Name:   n,
			Fields: fs,
		}
	}
	// Loop over the views again and this time fill the edges.
	for n, v := range vs {
		es := make(spec.Edges, len(v.Edges))
		for _, e := range v.Edges {
			es[e.Edge.Name] = spec.Edge{
				Schema: s.Components.Schemas[e.Name],
				Unique: e.Unique,
			}
		}
		s.Components.Schemas[n].Edges = es
	}
	return nil
}

// newField constructs a spec.Field out of a gen.Field.
func newField(f *gen.Field) (*spec.Field, error) {
	t, ok := oasTypes[f.Type.String()]
	if !ok {
		return nil, fmt.Errorf("no OAS-type exists for %q", f.Type.String())
	}
	e, err := exampleValue(f)
	if err != nil {
		return nil, err
	}
	return &spec.Field{
		Unique:   true,
		Required: !f.Optional,
		Type:     *t,
		Example:  e,
	}, nil
}

// exampleValue returns the user defined example value for this field.
func exampleValue(f *gen.Field) (interface{}, error) {
	var e interface{}
	a := Annotation{}
	if f.Annotations != nil && f.Annotations[a.Name()] != nil {
		if err := a.Decode(f.Annotations[a.Name()]); err != nil {
			return nil, err
		}
		e = a.Example
	}
	return e, nil
}

// requestBody returns the request-body to use for the given node and operation.
func requestBody(n *gen.Type, op string) (*spec.RequestBody, error) {
	req := &spec.RequestBody{}
	switch op {
	case createOperation:
		req.Description = fmt.Sprintf("%s to create", n.Name)
	case updateOperation:
		req.Description = fmt.Sprintf("%s properties to update", n.Name)
	default:
		return nil, fmt.Errorf("requestBody: unsupported operation %q", op)
	}
	fs := make(spec.Fields)
	for _, f := range n.Fields {
		if op == createOperation || !f.Immutable {
			sf, err := newField(f)
			if err != nil {
				return nil, err
			}
			fs[f.Name] = sf
		}
	}
	for _, e := range n.Edges {
		t, ok := oasTypes[e.Type.IDType.String()]
		if !ok {
			return nil, fmt.Errorf("no OAS-type exists for %q", e.Type.IDType.String())
		}
		fs[e.Name] = &spec.Field{
			Unique:   e.Unique,
			Required: !e.Optional,
			Type:     *t,
			Example:  nil, // TODO: Example for a unique / non-unique edge
		}
	}
	req.Content = spec.Content{
		spec.JSON: spec.MediaTypeObject{
			Schema: spec.Schema{
				Name:   fmt.Sprintf("%s%sRequest", n.Name, strcase.UpperCamelCase(op)),
				Fields: fs,
			},
		},
	}
	return req, nil
}

// paths adds all views to the specs schemas.
func paths(g *gen.Graph, s *spec.Spec) error {
	for _, n := range g.Nodes {
		ant := SchemaAnnotation{}
		if n.Annotations != nil && n.Annotations[ant.Name()] != nil {
			if err := ant.Decode(n.Annotations[ant.Name()]); err != nil {
				return err
			}
		}
		// Add schema operations.
		ops, err := nodeOperations(n)
		if err != nil {
			return err
		}
		// root for all operations on this node.
		root := "/" + rules.Pluralize(strcase.KebabCase(n.Name))
		// Create operation.
		if contains(ops, createOperation) {
			p := path(s, root)
			req, err := requestBody(n, createOperation)
			if err != nil {
				return err
			}
			v, err := newView(n, ant.CreateGroups)
			if err != nil {
				return err
			}
			rspName, err := v.Name()
			if err != nil {
				return err
			}
			p.Post = &spec.Operation{
				Summary:     fmt.Sprintf("creates a new %s", strcase.KebabCase(n.Name)),
				Description: fmt.Sprintf("Creates a new %s and persists it to storage", strcase.KebabCase(n.Name)),
				Tags:        []string{n.Name},
				OperationID: operationID(n, createOperation),
				RequestBody: req,
				Responses: map[string]spec.Response{
					strconv.Itoa(http.StatusOK): {
						Description: fmt.Sprintf("%s created", n.Name),
						Headers:     nil, // TODO
						Content: spec.Content{
							spec.JSON: spec.MediaTypeObject{
								Ref: s.Components.Schemas[rspName],
							},
						},
					},
				},
			}
		}
	}
	return nil
}

// path returns the correct spec.Path for the given root. Creates and sets a fresh instance if non does yet exist.
func path(s *spec.Spec, root string) *spec.Path {
	if s.Paths == nil {
		s.Paths = make(map[string]*spec.Path)
	}
	if _, ok := s.Paths[root]; !ok {
		s.Paths[root] = new(spec.Path)
	}
	return s.Paths[root]
}

// operationID generates a unique identifier for the given operation on the given node.
func operationID(n *gen.Type, op string) string {
	return op + n.Name
}

// // defaultString returns s if not "", d otherwise.
// func defaultString(s, d string) string {
// 	if s != "" {
// 		return s
// 	}
// 	return d
// }

func dump(out io.Writer, spec *spec.Spec) error {
	b, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	if out == nil {
		out, err = os.Create("openapi.json")
		if err != nil {
			return err
		}
		defer out.(*os.File).Close()
	}
	_, err = out.Write(b)
	return err
}
