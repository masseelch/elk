package elk

import (
	"encoding/json"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"errors"
	"github.com/masseelch/elk/policy"
	"github.com/masseelch/elk/spec"
	"io"
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
		specHooks      []Hook
	}
	// ExtensionOption allows managing Extension configuration using functional arguments.
	ExtensionOption func(*Extension) error
)

// NewExtension returns a new elk extension with default values.
func NewExtension(opts ...ExtensionOption) (*Extension, error) {
	ex := &Extension{
		config: &Config{HandlerPolicy: policy.Expose},
	}
	for _, opt := range opts {
		if err := opt(ex); err != nil {
			return nil, err
		}
	}
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

// EnableSpecGenerator enables the OpenAPI-Spec generator. Data will be written to given filename.
func EnableSpecGenerator(out string) ExtensionOption {
	return func(ex *Extension) error {
		if out == "" {
			return errors.New("spec filename cannot be empty")
		}
		ex.hooks = append(ex.hooks, ex.SpecGenerator(out))
		return nil
	}
}

// EnableHandlerGenerator enables generation of http crud handlers.
func EnableHandlerGenerator() ExtensionOption {
	return func(ex *Extension) error {
		ex.templates = []*gen.Template{HTTPTemplate}
		ex.easyjsonConfig = newEasyJsonConfig()
		ex.hooks = append(ex.hooks, EasyJSONGenerator(ex.easyjsonConfig))
		return nil
	}
}

// SpecHook registers the given Hook on the SpecGenerator.
func SpecHook(h Hook) ExtensionOption {
	return func(ex *Extension) error {
		if h == nil {
			return errors.New("hook cannot be nil")
		}
		ex.specHooks = append(ex.specHooks, h)
		return nil
	}
}

// SpecTitle sets the title of the Info block.
func SpecTitle(v string) ExtensionOption {
	return SpecHook(func(next Generator) Generator {
		return GenerateFunc(func(spec *spec.Spec) error {
			if err := next.Generate(spec); err != nil {
				return err
			}
			spec.Info.Title = v
			return nil
		})
	})
}

// SpecDescription sets the title of the Info block.
func SpecDescription(v string) ExtensionOption {
	return SpecHook(func(next Generator) Generator {
		return GenerateFunc(func(spec *spec.Spec) error {
			if err := next.Generate(spec); err != nil {
				return err
			}
			spec.Info.Description = v
			return nil
		})
	})
}

// SpecVersion sets the version of the Info block.
func SpecVersion(v string) ExtensionOption {
	return SpecHook(func(next Generator) Generator {
		return GenerateFunc(func(spec *spec.Spec) error {
			if err := next.Generate(spec); err != nil {
				return err
			}
			spec.Info.Version = v
			return nil
		})
	})
}

// TODO: Rest of Info block ...

// SpecDump dumps the current specs content to the given io.Writer.
func SpecDump(out io.Writer) ExtensionOption {
	return SpecHook(func(next Generator) Generator {
		return GenerateFunc(func(spec *spec.Spec) error {
			if err := next.Generate(spec); err != nil {
				return err
			}
			j, err := json.MarshalIndent(spec, "", "  ")
			if err != nil {
				return err
			}
			_, err = out.Write(j)
			return err
		})
	})
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
