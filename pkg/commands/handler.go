/*
Copyright © 2020 MasseElch <info@masseelch.de>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package commands

import (
	"bytes"
	"fmt"
	"github.com/facebook/ent/entc"
	"github.com/facebook/ent/entc/gen"
	"github.com/masseelch/elk/internal"
	"golang.org/x/tools/imports"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

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
		log.Fatal(err)
	}

	// Create the template
	tpl := template.New("handler").Funcs(gen.Funcs)
	for _, n := range []string{
		"header.tpl",
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

	var assets assets
	for _, n := range g.Nodes {
		assets.dirs = append(assets.dirs, filepath.Join(g.Config.Target, "handler"))
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

	return assets.format()
}

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
func (a assets) format() error {
	for _, file := range a.files {
		path := file.path
		src, err := imports.Process(path, file.content, nil)
		if err != nil {
			return fmt.Errorf("format file %s: %v", path, err)
		}
		if err := ioutil.WriteFile(path, src, 0644); err != nil {
			return fmt.Errorf("write file %s: %v", path, err)
		}
	}
	return nil
}
