package spec

const (
	// OpenAPI version 3.0.x is used.
	version = "3.0.3"

	JSON MediaType = "application/json"
)

type (
	Spec struct {
		Info         *Info            `json:"info"`
		Tags         []Tag            `json:"tags,omitempty"`
		Paths        map[string]*Path `json:"paths"`
		Components   Components       `json:"components"`
		Security     Security         `json:"security,omitempty"`
		ExternalDocs *ExternalDocs    `json:"externalDocs,omitempty"`
	}
	Tag struct {
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
	}
	Info struct {
		Title          string  `json:"title"`
		Description    string  `json:"description"`
		TermsOfService string  `json:"termsOfService"`
		Contact        Contact `json:"contact"`
		License        License `json:"license"`
		Version        string  `json:"version"`
	}
	Contact struct {
		Name  string `json:"name,omitempty"`
		Url   string `json:"url,omitempty"`
		Email string `json:"email,omitempty"`
	}
	License struct {
		Name string `json:"name"`
		URL  string `json:"url,omitempty"`
	}
	Path struct {
		Get        *Operation  `json:"get,omitempty"`
		Post       *Operation  `json:"post,omitempty"`
		Delete     *Operation  `json:"delete,omitempty"`
		Patch      *Operation  `json:"patch,omitempty"`
		Parameters []Parameter `json:"parameters,omitempty"`
	}
	Parameter struct {
		Name            string         `json:"name"`
		In              ParameterPlace `json:"in"`
		Description     string         `json:"description,omitempty"`
		Required        bool           `json:"required,omitempty"`
		Deprecated      bool           `json:"deprecated,omitempty"`
		AllowEmptyValue bool           `json:"allowEmptyValue,omitempty"`
		Schema          Type           `json:"schema"`
	}
	Operation struct {
		Summary      string                        `json:"summary,omitempty"`
		Description  string                        `json:"description,omitempty"`
		Tags         []string                      `json:"tags,omitempty"`
		ExternalDocs *ExternalDocs                 `json:"externalDocs,omitempty"`
		OperationID  string                        `json:"operationId"`
		Parameters   []*Parameter                  `json:"parameters,omitempty"`
		RequestBody  *RequestBody                  `json:"requestBody,omitempty"`
		Responses    map[string]*OperationResponse `json:"responses"`
		Deprecated   bool                          `json:"deprecated,omitempty"`
		Security     Security                      `json:"security,omitempty"`
	}
	Security          []map[string][]string
	OperationResponse struct {
		Ref      *Response
		Response Response
	}
	ExternalDocs struct {
		Description string `json:"description"`
		URL         string `json:"url"`
	}
	RequestBody struct {
		Description string  `json:"description"`
		Content     Content `json:"content"`
	}
	Content         map[MediaType]MediaTypeObject
	MediaType       string
	MediaTypeObject struct {
		Unique bool    `json:"-"`
		Ref    *Schema `json:"-"`
		Schema Schema  `json:"schema"`
	}
	Schema struct {
		Name   string
		Fields Fields
		Edges  Edges
	}
	Fields map[string]*Field
	Field  struct {
		Type
		Unique   bool        `json:"-"`
		Required bool        `json:"-"`
		Example  interface{} `json:"example,omitempty"`
	}
	Type struct {
		Type   string `json:"type,omitempty"`
		Format string `json:"format,omitempty"`
		Items  *Type  `json:"items,omitempty"`
	}
	Edges map[string]Edge
	Edge  struct {
		Schema Schema  `json:"schema"`
		Ref    *Schema `json:"-"`
		Unique bool    `json:"-"`
	}
	Response struct {
		Name        string               `json:"-"`
		Description string               `json:"description"`
		Headers     map[string]Parameter `json:"headers,omitempty"`
		Content     *Content             `json:"content,omitempty"`
	}
	Components struct {
		Schemas         map[string]*Schema        `json:"schemas"`
		Responses       map[string]*Response      `json:"responses"`
		Parameters      map[string]Parameter      `json:"parameters"`
		SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty"`
	}
	SecurityScheme struct {
		Type             string      `json:"type"`
		Description      string      `json:"description,omitempty"`
		Name             string      `json:"name,omitempty"`
		In               string      `json:"in,omitempty"`
		Scheme           string      `json:"scheme,omitempty"`
		BearerFormat     string      `json:"bearerFormat,omitempty"`
		Flows            *OAuthFlows `json:"flows,omitempty"`
		OpenIdConnectUrl string      `json:"openIdConnectUrl,omitempty"`
	}
	OAuthFlows struct {
		Implicit          *OAuthFlow `json:"implicit,omitempty"`
		Password          *OAuthFlow `json:"password,omitempty"`
		ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty"`
		AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty"`
	}
	OAuthFlow struct {
		AuthorizationUrl string            `json:"authorizationUrl,omitempty"`
		TokenUrl         string            `json:"tokenUrl,omitempty"`
		RefreshUrl       string            `json:"refreshUrl,omitempty"`
		Scopes           map[string]string `json:"scopes,omitempty"`
	}
)
