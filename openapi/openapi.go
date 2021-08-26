package openapi

import "net/url"

// OpenAPI spec 3.0.x is used.
const openApiVersion = "3.0.3"

type (
	Spec struct {
		Info    *Info
		Servers []Server
		Paths   []Path
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
	Server struct {
		URL         url.URL
		Description string
		Variables   []ServerVariable
	}
	ServerVariable struct {
		Enum        []string
		Default     string
		Description string
	}
	Path struct {
		Summary     string
		Description string
		Method      string
		Operation   Operation
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
	}
	ExternalDocs struct {
		Description string
		URL         url.URL
	}
	RequestBody struct {
		Description string
		Content     map[MediaType]MediaTypeObject
	}
	MediaTypeObject struct {
		Schema  *Schema
		Example interface{}
	}
	Schema struct {
	}
	// SpecOption allows managing OpenAPI configuration using functional arguments.
	SpecOption func(*Spec) error
)

func New(opts ...SpecOption) (*Spec, error) {
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
func Title(t string) SpecOption {
	return func(spec *Spec) error {
		spec.getInfo().Title = t
		return nil
	}
}

// Description sets the title of the Info block.
func Description(d string) SpecOption {
	return func(spec *Spec) error {
		spec.getInfo().Description = d
		return nil
	}
}

// Version sets the title of the Info block.
func Version(d string) SpecOption {
	return func(spec *Spec) error {
		spec.getInfo().Version = d
		return nil
	}
}

// And so on ...

func (spec Spec) getInfo() *Info {
	if spec.Info == nil {
		spec.Info = &Info{}
	}
	return spec.Info
}
