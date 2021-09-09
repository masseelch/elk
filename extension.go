package elk

import (
	"encoding/json"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"errors"
	"github.com/masseelch/elk/spec"
	"io"
)

type (
	Config struct {
		// HandlerPolicy defines the default policy for handler generation.
		// It is used if no policy is set on a (sub-)resource.
		// Defaults to policy.Expose.
		HandlerPolicy Policy
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
	// HandlerOption allows managing RESTGenerator configuration using function arguments.
	HandlerOption ExtensionOption
)

// NewExtension returns a new elk extension with default values.
func NewExtension(opts ...ExtensionOption) (*Extension, error) {
	ex := &Extension{
		config: &Config{HandlerPolicy: Expose},
	}
	for _, opt := range opts {
		if err := opt(ex); err != nil {
			return nil, err
		}
	}
	if len(ex.hooks) == 0 {
		return nil, errors.New(`no generator enabled: enable one by providing either "EnableSpecGenerator()" or "EnableHandlerGenerator()" to "NewExtension()"`)
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
func DefaultHandlerPolicy(p Policy) ExtensionOption {
	return func(ex *Extension) error {
		if err := p.Validate(); err != nil {
			return err
		}
		ex.config.HandlerPolicy = p
		return nil
	}
}

// HandlerEasyJsonConfig sets a custom EasyJsonConfig.
func HandlerEasyJsonConfig(c EasyJsonConfig) HandlerOption {
	return func(ex *Extension) error {
		ex.easyjsonConfig = c
		return nil
	}
}

// GenerateSpec enables the OpenAPI-Spec generator. Data will be written to given filename.
func GenerateSpec(out string, hooks ...Hook) ExtensionOption {
	return func(ex *Extension) error {
		if out == "" {
			return errors.New("spec filename cannot be empty")
		}
		ex.hooks = append(ex.hooks, ex.SpecGenerator(out))
		if len(hooks) > 0 {
			ex.specHooks = append(ex.specHooks, hooks...)
		}
		return nil
	}
}

// GenerateHandlers enables generation of http crud handlers.
func GenerateHandlers(opts ...HandlerOption) ExtensionOption {
	return func(ex *Extension) error {
		ex.templates = []*gen.Template{HTTPTemplate}
		ex.easyjsonConfig = newEasyJsonConfig()
		ex.hooks = append(ex.hooks, EasyJSONGenerator(ex.easyjsonConfig))
		for _, opt := range opts {
			if err := opt(ex); err != nil {
				return err
			}
		}
		return nil
	}
}

// SpecTitle sets the title of the Info block.
func SpecTitle(v string) Hook {
	return func(next Generator) Generator {
		return GenerateFunc(func(spec *spec.Spec) error {
			if err := next.Generate(spec); err != nil {
				return err
			}
			spec.Info.Title = v
			return nil
		})
	}
}

// SpecDescription sets the title of the Info block.
func SpecDescription(v string) Hook {
	return func(next Generator) Generator {
		return GenerateFunc(func(spec *spec.Spec) error {
			if err := next.Generate(spec); err != nil {
				return err
			}
			spec.Info.Description = v
			return nil
		})
	}
}

// SpecVersion sets the version of the Info block.
func SpecVersion(v string) Hook {
	return func(next Generator) Generator {
		return GenerateFunc(func(spec *spec.Spec) error {
			if err := next.Generate(spec); err != nil {
				return err
			}
			spec.Info.Version = v
			return nil
		})
	}
}

// TODO: Rest of Info block ...

// SpecSecuritySchemes sets the security schemes of the Components block.
func SpecSecuritySchemes(schemes map[string]spec.SecurityScheme) Hook {
	return func(next Generator) Generator {
		return GenerateFunc(func(spec *spec.Spec) error {
			if err := next.Generate(spec); err != nil {
				return err
			}
			spec.Components.SecuritySchemes = schemes
			return nil
		})
	}
}

// SpecSecurity sets the global security Spec.
func SpecSecurity(sec spec.Security) Hook {
	return func(next Generator) Generator {
		return GenerateFunc(func(spec *spec.Spec) error {
			if err := next.Generate(spec); err != nil {
				return err
			}
			spec.Security = sec
			return nil
		})
	}
}

// SpecDump dumps the current specs content to the given io.Writer.
func SpecDump(out io.Writer) Hook {
	return func(next Generator) Generator {
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
