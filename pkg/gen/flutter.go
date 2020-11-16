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
		"flutter/client.tpl",
		"flutter/client_provider.tpl",
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
			filepath.Join(g.Config.Target, "client"),
		},
	}

	for _, n := range g.Nodes {
		b := bytes.NewBuffer(nil)
		if err := tpl.ExecuteTemplate(b, "model", n); err != nil {
			panic(err)
		}
		assets.files = append(assets.files, file{
			path:    filepath.Join(g.Config.Target, "model", fmt.Sprintf("%s.dart", gen.Funcs["snake"].(func(string) string)(n.Name))),
			content: b.Bytes(),
		})

		// Only generate the client if the generation should not be skipped.
		if n.Annotations["HandlerGen"] == nil || !n.Annotations["HandlerGen"].(map[string]interface{})["Skip"].(bool) {
			b = bytes.NewBuffer(nil)
			if err := tpl.ExecuteTemplate(b, "client", n); err != nil {
				panic(err)
			}
			assets.files = append(assets.files, file{
				path:    filepath.Join(g.Config.Target, "client", fmt.Sprintf("%s.dart", gen.Funcs["snake"].(func(string) string)(n.Name))),
				content: b.Bytes(),
			})
		}
	}

	b := bytes.NewBuffer(nil)
	if err := tpl.ExecuteTemplate(b, "client/provider", g); err != nil {
		return err
	}
	assets.files = append(assets.files, file{
		path:    filepath.Join(g.Config.Target, "client", "provider.dart"),
		content: b.Bytes(),
	})

	if err := assets.write(); err != nil {
		return err
	}

	return assets.formatDart()
}
