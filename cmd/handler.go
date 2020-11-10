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
package cmd

import (
	"bytes"
	"fmt"
	"github.com/facebook/ent/entc"
	"github.com/facebook/ent/entc/gen"
	"github.com/masseelch/elk/internal"
	"github.com/spf13/cobra"
	"golang.org/x/tools/imports"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

// handlerCmd represents the handler command
var handlerCmd = &cobra.Command{
	Use:   "handler",
	Short: "generate api handlers for your defined schemas",
	Run: func(cmd *cobra.Command, args []string) {
		// Flags
		s, err := cmd.Flags().GetString("source")
		if err != nil {
			log.Fatal(err)
		}
		t, err := cmd.Flags().GetString("target")
		if err != nil {
			log.Fatal(err)
		}

		// Load the graph
		g, err := entc.LoadGraph(s, &gen.Config{
			Target: t,
		})
		if err != nil {
			log.Fatal(err)
		}

		// Create the template
		tpl := template.New("handler").Funcs(gen.Funcs)
		tpl = template.Must(tpl.Parse(string(internal.MustAsset("header.tpl"))))
		tpl = template.Must(tpl.Parse(string(internal.MustAsset("handler/handler.tpl"))))
		tpl = template.Must(tpl.Parse(string(internal.MustAsset("handler/create.tpl"))))
		tpl = template.Must(tpl.Parse(string(internal.MustAsset("handler/read.tpl"))))
		// tpl = template.Must(tpl.Parse(string(internal.MustAsset("handler/update.tpl"))))
		// tpl = template.Must(tpl.Parse(string(internal.MustAsset("handler/delete.tpl"))))
		tpl = template.Must(tpl.Parse(string(internal.MustAsset("handler/list.tpl"))))

		var assets assets
		for _, n := range g.Nodes {
			// assets.dirs = append(assets.dirs, filepath.Join(g.Config.Target, "handler"))
			b := bytes.NewBuffer(nil)
			if err := tpl.Execute(b, n); err != nil {
				panic(err)
			}
			assets.files = append(assets.files, file{
				path:    filepath.Join(g.Config.Target, fmt.Sprintf("%s.go", gen.Funcs["snake"].(func(string) string)(n.Name))),
				content: b.Bytes(),
			})

			if err := assets.write(); err != nil {
				panic(err)
			}

			if err := assets.format(); err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	generateCmd.AddCommand(handlerCmd)
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
