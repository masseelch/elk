package spec

import "net/url"

const (
	// OpenAPI spec 3.0.x is used.
	openApiVersion = "3.0.3"

	JSON MediaType = "application/json"
)

// OASTyper is the interface a custom non-primitive GoType has to implement to be convertable to an OAS-type.
type OASTyper interface {
	OASType() Type
}

type (
	Spec struct {
		Info       *Info      `json:"info"`
		Tags       []Tag      `json:"tags"`
		Paths      []Path     `json:"paths"`
		Components Components `json:"components"`
	}
	Tag struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	Info struct {
		Title          string  `json:"title"`
		Description    string  `json:"description"`
		TermsOfService string  `json:"terms_of_service"`
		Contact        Contact `json:"contact"`
		License        License `json:"license"`
		Version        string  `json:"version"`
	}
	URL     url.URL
	Contact struct {
		Name  string `json:"name"`
		Url   URL    `json:"url"`
		Email string `json:"email"`
	}
	License struct {
		Name string `json:"name"`
		URL  URL    `json:"url"`
	}
	Path struct {
		Summary     string      `json:"summary"`
		Description string      `json:"description"`
		Get         Operation   `json:"get"`
		Post        Operation   `json:"post"`
		Delete      Operation   `json:"delete"`
		Patch       Operation   `json:"patch"`
		Parameters  []Parameter `json:"parameters"`
	}
	Parameter struct {
		Name            string         `json:"name"`
		In              ParameterPlace `json:"in"`
		Description     string         `json:"description"`
		Required        bool           `json:"required"`
		Deprecated      bool           `json:"deprecated"`
		AllowEmptyValue bool           `json:"allow_empty_value"`
	}
	Operation struct {
		Tags         []string     `json:"tags"`
		Summary      string       `json:"summary"`
		Description  string       `json:"description"`
		ExternalDocs ExternalDocs `json:"external_docs"`
		OperationID  string       `json:"operation_id"`
		Parameters   []Parameter  `json:"parameters"`
		RequestBody  RequestBody  `json:"request_body"`
		Responses    []Response   `json:"responses"`
		Deprecated   bool         `json:"deprecated"`
	}
	ExternalDocs struct {
		Description string  `json:"description"`
		URL         url.URL `json:"url"`
	}
	RequestBody struct {
		Description string  `json:"description"`
		Content     Content `json:"content"`
	}
	Content         map[MediaType]MediaTypeObject
	MediaType       string
	MediaTypeObject struct {
		Schema  Schema      `json:"schema"`
		Example interface{} `json:"example"`
	}
	Response struct {
		Code        int                  `json:"code"`
		Description string               `json:"description"`
		Headers     map[string]Parameter `json:"headers"`
		Content     Content              `json:"content"`
	}
	Components struct {
		Schemas    map[string]Schema    `json:"schemas"`
		Responses  map[string]Response  `json:"responses"`
		Parameters map[string]Parameter `json:"parameters"`
		// ... TODO
	}
	Fields map[string]Field
	Edges  map[string]Edge
	Schema struct {
		Fields Fields `json:"fields"`
		Edges  Edges  `json:"edges"`
	}
	// SchemaRef is used to serialize a #ref instead of an object.
	SchemaRef Schema
	Field     struct {
		Type
		Required bool        `json:"-"`
		Example  interface{} `json:"example,omitempty"`
	}
	Edge struct {
	}
	Property struct {
		Type
		Example interface{} `json:"example"`
	}
	Type struct {
		Type   string `json:"type"`
		Format string `json:"format,omitempty"`
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
	// if spec.Components == nil {
	// 	spec.Components = &Components{}
	// }
	c := &spec.Components
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
