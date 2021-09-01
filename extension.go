package elk

import (
	"encoding/json"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/masseelch/elk/policy"
	"github.com/masseelch/elk/spec"
)

type (
	Config struct {
		// HandlerPolicy defines the default policy for handler generation.
		// It is used if no policy is set on a (sub-)resource.
		// Defaults to policy.Expose.
		HandlerPolicy policy.Policy
	}
	// Extension implements entc.Extension interface for providing http handler code generation.
	Extension struct {
		entc.DefaultExtension
		easyjsonConfig EasyJsonConfig
		hooks          []gen.Hook
		templates      []*gen.Template
		config         *Config
		// If non-nil the generator generates an OpenAPI-Specification for the defined schemas.
		spec *spec.Spec
	}
	// ExtensionOption allows managing Extension configuration using functional arguments.
	ExtensionOption func(*Extension) error
)

// NewExtension returns a new elk extension with default values.
func NewExtension(opts ...ExtensionOption) (*Extension, error) {
	ex := &Extension{
		templates:      []*gen.Template{HTTPTemplate},
		easyjsonConfig: newEasyJsonConfig(),
		config: &Config{
			HandlerPolicy: policy.Expose,
		},
	}
	for _, opt := range opts {
		if err := opt(ex); err != nil {
			return nil, err
		}
	}
	ex.hooks = append(ex.hooks, EasyJSONGenerator(ex.easyjsonConfig))
	return ex, nil
}

// Templates of the Extension.
func (e *Extension) Templates() []*gen.Template {
	return e.templates
}

// Hooks of the Extension.
func (e *Extension) Hooks() []gen.Hook {
	return e.hooks
}

// Annotations of the Extension.
func (e *Extension) Annotations() []entc.Annotation {
	return []entc.Annotation{e.config}
}

// DefaultHandlerPolicy sets the policy.Policy to use of none is given on a (sub-)schema.
func DefaultHandlerPolicy(p policy.Policy) ExtensionOption {
	return func(ex *Extension) error {
		if err := p.Validate(); err != nil {
			return err
		}
		ex.config.HandlerPolicy = p
		return nil
	}
}

// WithEasyJsonConfig sets a custom EasyJsonConfig.
func WithEasyJsonConfig(c EasyJsonConfig) ExtensionOption {
	return func(ex *Extension) error {
		ex.easyjsonConfig = c
		return nil
	}
}

// WithOpenAPISpec enables the OpenAPI-Spec generator, which will merge into the given spec.
func WithOpenAPISpec(spec *spec.Spec) ExtensionOption {
	return func(ex *Extension) error {
		ex.spec = spec
		ex.hooks = append(ex.hooks, SpecGenerator(spec))
		return nil
	}
}

// Name implements entc.Annotation interface.
func (c Config) Name() string {
	return "ElkConfig"
}

// Decode from ent.
func (c *Config) Decode(o interface{}) error {
	buf, err := json.Marshal(o)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, c)
}

var _ entc.Annotation = (*Config)(nil)
