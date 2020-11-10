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

func Flutter(source string, target string) error {
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
	tpl := template.New("flutter").Funcs(gen.Funcs).Funcs(template.FuncMap{"dartType": dartType})
	for _, n := range []string{
		"header/dart.tpl",
		"flutter/model.tpl",
		"flutter/repository.tpl",
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

	assets := assets{
		dirs: []string{
			filepath.Join(g.Config.Target, "model"),
			filepath.Join(g.Config.Target, "repository"),
		},
	}

	for _, n := range g.Nodes {
		m := bytes.NewBuffer(nil)
		if err := tpl.ExecuteTemplate(m, "model", n); err != nil {
			panic(err)
		}
		assets.files = append(assets.files, file{
			path:    filepath.Join(g.Config.Target, "model", fmt.Sprintf("%s.dart", gen.Funcs["snake"].(func(string) string)(n.Name))),
			content: m.Bytes(),
		})

		r := bytes.NewBuffer(nil)
		if err := tpl.ExecuteTemplate(r, "repository", n); err != nil {
			panic(err)
		}
		assets.files = append(assets.files, file{
			path:    filepath.Join(g.Config.Target, "repository", fmt.Sprintf("%s.dart", gen.Funcs["snake"].(func(string) string)(n.Name))),
			content: r.Bytes(),
		})
	}

	if err := assets.write(); err != nil {
		return err
	}

	return assets.formatDart()
}
