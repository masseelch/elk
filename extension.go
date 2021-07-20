package elk

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

type (
	// Extension implements entc.Extension interface for providing http handler code generation.
	Extension struct {
		entc.DefaultExtension
		// cfg       Config
		hooks     []gen.Hook
		templates []*gen.Template
	}
	// ExtensionOption allows to manage Extension configuration using functional arguments.
	ExtensionOption func(*Extension) error
)

func NewExtension(opts ...ExtensionOption) (*Extension, error) {
	ex := &Extension{
		templates: HTTPTemplates,
		hooks:     []gen.Hook{AddGroupsTag},
	}
	for _, opt := range opts {
		if err := opt(ex); err != nil {
			return nil, err
		}
	}
	return ex, nil
}

func (e *Extension) Templates() []*gen.Template {
	return e.templates
}

func (e *Extension) Hooks() []gen.Hook {
	return e.hooks
}
