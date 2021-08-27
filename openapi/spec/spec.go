package spec

import "net/url"

const (
	// OpenAPI spec 3.0.x is used.
	openApiVersion = "3.0.3"

	JSON MediaType = "application/json"
)

type (
	Spec struct {
		Info       *Info
		Tags       []Tag
		Paths      []Path
		Components Components
	}
	Tag struct {
		Name        string
		Description string
	}
	Info struct {
		Title          string
		Description    string
		TermsOfService string
		Contact        *Contact
		License        *License
		Version        string
	}
	Contact struct {
		Name  string
		Url   url.URL
		Email string
	}
	License struct {
		Name string
		URL  url.URL
	}
	Path struct {
		Summary     string
		Description string
		Get         Operation
		Post        Operation
		Delete      Operation
		Patch       Operation
		Parameters  []Parameter
	}
	Parameter struct {
		Name            string
		In              ParameterPlace
		Description     string
		Required        bool
		Deprecated      bool
		AllowEmptyValue bool
	}
	Operation struct {
		Tags         []string
		Summary      string
		Description  string
		ExternalDocs ExternalDocs
		OperationID  string
		Parameters   []Parameter
		RequestBody  RequestBody
		Responses    []Response
		Deprecated   bool
	}
	ExternalDocs struct {
		Description string
		URL         url.URL
	}
	RequestBody struct {
		Description string
		Content     Content
	}
	Content         map[MediaType]MediaTypeObject
	MediaType       string
	MediaTypeObject struct {
		Schema  Schema
		Example interface{}
	}
	Response struct {
		Code        int
		Description string
		Headers     map[string]Parameter
		Content     Content
	}
	Components struct {
		Schemas    map[string]Schema
		Responses  map[string]Response
		Parameters map[string]Parameter
		// ... TODO
	}
	Schema struct {
		Fields map[string]Field
		Edges  map[string]Edge
	}
	Field struct {
		Required bool
		Type     string
		Format   string
		Example  interface{}
	}
	Edge struct {
	}
	Property struct {
		Type    string
		Format  string
		Example interface{}
	}
	// Option allows managing spec-configuration using functional arguments.
	Option func(*Spec) error
)

func New(opts ...Option) (*Spec, error) {
	oa := &Spec{
		Info: &Info{
			Title:       "Ent Schema API",
			Description: "This is a auto generated API description made out of an Ent schema definition",
			Version:     "0.0.0",
		},
	}
	for _, opt := range opts {
		if err := opt(oa); err != nil {
			return nil, err
		}
	}
	return oa, nil
}

// Title sets the title of the Info block.
func Title(v string) Option {
	return func(spec *Spec) error {
		spec.getInfo().Title = v
		return nil
	}
}

// Description sets the title of the Info block.
func Description(v string) Option {
	return func(spec *Spec) error {
		spec.getInfo().Description = v
		return nil
	}
}

// Version sets the title of the Info block.
func Version(v string) Option {
	return func(spec *Spec) error {
		spec.getInfo().Version = v
		return nil
	}
}

// Tags sets the Tags block.
func Tags(ts ...Tag) Option {
	return func(spec *Spec) error {
		spec.Tags = ts
		return nil
	}
}

// And so on ...

func (spec *Spec) WarmUp() {
	// Ensure a non nil Info block.
	_ = spec.getInfo()
	// Ensure non nil maps in Components block.
	c := spec.Components
	if c.Schemas == nil {
		c.Schemas = make(map[string]Schema)
	}
	if c.Responses == nil {
		c.Responses = make(map[string]Response)
	}
	if c.Parameters == nil {
		c.Parameters = make(map[string]Parameter)
	}
}

func (spec Spec) getInfo() *Info {
	if spec.Info == nil {
		spec.Info = &Info{}
	}
	return spec.Info
}
