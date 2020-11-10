package gen

import (
	"bytes"
	"fmt"
	"github.com/facebook/ent/entc"
	"github.com/facebook/ent/entc/gen"
	"github.com/masseelch/elk/internal"
	"path/filepath"
	"text/template"
)

func Handler(source string, target string) error {
	cfg := &gen.Config{Target: target}
	if cfg.Target == "" {
		abs, err := filepath.Abs(source)
		if err != nil {
			return err
		}
		// Default target-path for codegen is one dir above the schema.
		cfg.Target = filepath.Dir(abs)
	}

	// Load the graph
	g, err := entc.LoadGraph(source, cfg)
	if err != nil {
		return err
	}

	// Create the template
	tpl := template.New("handler").Funcs(gen.Funcs)
	for _, n := range []string{
		"header/go.tpl",
		"handler/handler.tpl",
		"handler/create.tpl",
		"handler/read.tpl",
		"handler/update.tpl",
		// "handler/delete.tpl",
		"handler/list.tpl",
	} {
		d, err := internal.Asset(n)
		if err != nil {
			return err
		}
		tpl, err = tpl.Parse(string(d))
		if err != nil {
			return err
		}
	}

	assets := assets{dirs: []string{filepath.Join(g.Config.Target, "handler")}}
	for _, n := range g.Nodes {
		b := bytes.NewBuffer(nil)
		if err := tpl.Execute(b, n); err != nil {
			panic(err)
		}
		assets.files = append(assets.files, file{
			path:    filepath.Join(g.Config.Target, "handler", fmt.Sprintf("%s.go", gen.Funcs["snake"].(func(string) string)(n.Name))),
			content: b.Bytes(),
		})

	}

	if err := assets.write(); err != nil {
		return err
	}

	return assets.formatGo()
}
