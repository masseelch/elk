package elk

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

type (
	// Extension implements entc.Extension interface for providing http handler code generation.
	Extension struct {
		entc.DefaultExtension
		easyjsonConfig EasyJsonConfig
		hooks          []gen.Hook
		templates      []*gen.Template
	}
	// ExtensionOption allows to manage Extension configuration using functional arguments.
	ExtensionOption func(*Extension) error
)

// NewExtension returns a new elk extension with default values.
func NewExtension(opts ...ExtensionOption) (*Extension, error) {
	ex := &Extension{
		templates:      HTTPTemplates,
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

func (e *Extension) Templates() []*gen.Template {
	return e.templates
}

func (e *Extension) Hooks() []gen.Hook {
	return e.hooks
}

func WithEasyJsonConfig(c EasyJsonConfig) ExtensionOption {
	return func(ex *Extension) error {
		ex.easyjsonConfig = c
		return nil
	}
}
