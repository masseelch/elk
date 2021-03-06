package gen

import (
	"bytes"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/masseelch/elk/internal"
	"path/filepath"
	"text/template"
)

func Handler(c *Config) error {
	cfg := &gen.Config{Target: c.Target}
	if cfg.Target == "" {
		abs, err := filepath.Abs(c.Source)
		if err != nil {
			return err
		}
		// Default target-path for codegen is one dir above the schema.
		cfg.Target = filepath.Dir(abs)
	}

	// Load the graph
	g, err := entc.LoadGraph(c.Source, cfg)
	if err != nil {
		return err
	}

	// Create the template
	tpl := template.New("handler").Funcs(gen.Funcs).Funcs(template.FuncMap{
		"eagerLoadedEdges": eagerLoadedEdges,
		"eagerLoadBuilder": eagerLoadQueryByGroupKeyRecursive,
		"pkgImports":       pkgImports,
		"lowerFirst":       lowerFirst,
		"visited":          visited,
	})

	// Attach header template.
	tpl, err = tpl.Parse(string(internal.MustAsset("header/go.tpl")))
	if err != nil {
		return err
	}

	// Load all handler templates.
	ts, err := internal.AssetDir("handler")
	if err != nil {
		return err
	}
	for _, n := range ts {
		d, err := internal.Asset("handler/" + n)
		if err != nil {
			return err
		}
		tpl, err = tpl.Parse(string(d))
		if err != nil {
			return err
		}
	}

	// Generate the code.
	assets := assets{dirs: []string{filepath.Join(g.Config.Target, "handler")}}
	b := bytes.NewBuffer(nil)
	if err := tpl.Execute(b, g); err != nil {
		panic(err)
	}
	assets.files = append(assets.files, file{
		path:    filepath.Join(g.Config.Target, "handler", "handler.go"),
		content: b.Bytes(),
	})

	// Write and format the generated code files.
	if err := assets.write(); err != nil {
		return err
	}

	return assets.formatGo()
}
