package elk

import (
	"entgo.io/contrib/entoas"
	"entgo.io/ent/entc/gen"
)

type (
	// Extension implements entc.Extension interface for providing http handler code generation.
	Extension struct {
		*entoas.Extension

		// easyjsonConfig EasyJsonConfig
		oasOptions []entoas.ExtensionOption
	}
	// ExtensionOption allows managing Extension configuration using functional arguments.
	ExtensionOption func(*Extension) error
)

// NewExtension returns a new elk extension with default values.
func NewExtension(opts ...ExtensionOption) (*Extension, error) {
	ex := &Extension{}
	oas, err := entoas.NewExtension()
	if err != nil {
		return nil, err
	}
	ex.Extension = oas
	for _, opt := range opts {
		if err := opt(ex); err != nil {
			return nil, err
		}
	}
	return ex, nil
}

// EntoasOptions takes options normally passed to the entoas extension and passes them on to the wrapped extension.
func EntoasOptions(opts ...entoas.ExtensionOption) ExtensionOption {
	return func(ex *Extension) error {
		for _, opt := range opts {
			if err := opt(ex.Extension); err != nil {
				return err
			}
		}
		return nil
	}
}

// // HandlerEasyJsonConfig sets a custom EasyJsonConfig.
// func HandlerEasyJsonConfig(c EasyJsonConfig) HandlerOption {
// 	return func(ex *Extension) error {
// 		ex.easyjsonConfig = c
// 		return nil
// 	}
// }

// Hooks of the Extension.
func (ex *Extension) Hooks() []gen.Hook {
	return append(
		ex.Extension.Hooks(),
		// TODO: easyJson
	)
}

// Templates of the Extension.
func (ex *Extension) Templates() []*gen.Template {
	return []*gen.Template{HTTPTemplates}
}
