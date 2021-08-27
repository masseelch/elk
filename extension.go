package elk

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/masseelch/elk/openapi"
	"github.com/masseelch/elk/openapi/spec"
)

type (
	// Extension implements entc.Extension interface for providing http handler code generation.
	Extension struct {
		entc.DefaultExtension
		easyjsonConfig EasyJsonConfig
		hooks          []gen.Hook
		templates      []*gen.Template
		// If non-nil the generator will generate an OpenAPI-Specification for the defined schemas.
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
	}
	for _, opt := range opts {
		if err := opt(ex); err != nil {
			return nil, err
		}
	}
	ex.hooks = append(ex.hooks, GenerateEasyJSON(ex.easyjsonConfig))
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
		ex.hooks = append(ex.hooks, openapi.Generator(spec))
		return nil
	}
}
