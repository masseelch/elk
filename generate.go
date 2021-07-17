package elk

import (
	"bytes"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"errors"
	"fmt"
	"github.com/masseelch/elk/internal"
	"golang.org/x/tools/imports"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

func Flutter(schemaPath, targetPath string) error {
	if targetPath == "" {
		abs, err := filepath.Abs(schemaPath)
		if err != nil {
			return err
		}
		// Default target-path for codegen is one dir above the schema.
		targetPath = filepath.Dir(abs)
	}

	// Load the graph.
	g, err := entc.LoadGraph(schemaPath, &gen.Config{})
	if err != nil {
		return err
	}

	// Parse templates.
	t := template.New("").Funcs(gen.Funcs).Funcs(TemplateFuncs)
	tpls, err := internal.AssetDir("template/flutter")
	if err != nil {
		return err
	}
	for _, tpl := range tpls {
		_, err = t.Parse(string(internal.MustAsset("template/flutter/" + tpl)))
		if err != nil {
			return err
		}
	}

	// Run the templates.
	assets := assets{
		dirs: []string{
			filepath.Join(targetPath, "model"),
			filepath.Join(targetPath, "client"),
		},
	}
	b := bytes.NewBuffer(nil)
	for _, n := range g.Nodes {
		// Generate model for node.
		if err := t.ExecuteTemplate(b, "flutter/model", n); err != nil {
			return fmt.Errorf("execute template %q: %w", "flutter/model", err)
		}
		assets.files = append(assets.files, file{
			path:    filepath.Join(targetPath, "model", fmt.Sprintf("%s.dart", n.Label())),
			content: b.Bytes(),
		})
		b.Reset()

		// Generate client for node.
		if err := t.ExecuteTemplate(b, "flutter/client", n); err != nil {
			return fmt.Errorf("execute template %q: %w", "flutter/client", err)
		}
		assets.files = append(assets.files, file{
			path:    filepath.Join(targetPath, "client", fmt.Sprintf("%s.dart", n.Label())),
			content: b.Bytes(),
		})
		b.Reset()
	}

	// Write and format assets only if template execution
	// finished successfully.
	if err := assets.write(); err != nil {
		return err
	}

	// "dartfmt" can only run on file system.
	return assets.formatDart()
}

type (
	file struct {
		path    string
		content []byte
	}
	assets struct {
		dirs  []string
		files []file
	}
)

// write files and dirs in the assets.
func (a assets) write() error {
	for _, dir := range a.dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("create dir %q: %w", dir, err)
		}
	}
	for _, file := range a.files {
		if err := ioutil.WriteFile(file.path, file.content, 0644); err != nil {
			return fmt.Errorf("write file %q: %w", file.path, err)
		}
	}
	return nil
}

// format runs "goimports" on all assets.
func (a assets) formatGo() error {
	for _, file := range a.files {
		path := file.path
		src, err := imports.Process(path, file.content, nil)
		if err != nil {
			return fmt.Errorf("format file %s: %w", path, err)
		}
		if err := ioutil.WriteFile(path, src, 0644); err != nil {
			return fmt.Errorf("write file %s: %w", path, err)
		}
	}
	return nil
}

// formatDart runs "dartfmt" on all assets.
func (a assets) formatDart() error {
	args := append([]string{"-w", "--line-length=120"}, a.dirs...)
	_, err := exec.Command("dartfmt", args...).Output()
	if err != nil {
		if err == exec.ErrNotFound {
			// "dartfmt" is not available
			return errors.New("dartfmt: command not found")
		} else {
			if err, ok := err.(*exec.ExitError); ok {
				return fmt.Errorf("dartfmt: %w, %s", err, err.Stderr)
			}

			return fmt.Errorf("dartfmt: %w", err)
		}
	}

	return nil
}
