package gen

import (
	"fmt"
	"github.com/facebook/ent/entc/gen"
	"github.com/facebook/ent/schema/field"
	"golang.org/x/tools/imports"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

var (
	dartTypeNames = map[string]string{
		"invalid":   "dynamic",
		"bool":      "bool",
		"time.Time": "DateTime",
		// "JSON":    "Map<String, dynamic>",
		// "UUID":    "String",
		// "bytes":   "dynamic",
		"enum":    "String",
		"string":  "String",
		"int":     "int",
		"int8":    "int",
		"int16":   "int",
		"int32":   "int",
		"int64":   "int",
		"uint":    "int",
		"uint8":   "int",
		"uint16":  "int",
		"uint32":  "int",
		"uint64":  "int",
		"float32": "double",
		"float64": "double",
	}
)

type (
	dartField struct {
		Type      string
		Converter string
		Field     *gen.Field
		Edge      *gen.Edge
	}
	file struct {
		path    string
		content []byte
	}
	assets struct {
		dirs  []string
		files []file
	}
	Config struct {
		Source string
		Target string
	}
)

func (d dartField) IsEdge() bool {
	return d.Edge != nil
}

func (d dartField) Name() string {
	if d.IsEdge() {
		return d.Edge.Name
	}

	return d.Field.Name
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

// formatGo runs "goimports" on all assets.
func (a assets) formatGo() error {
	for _, file := range a.files {
		path := file.path
		src, err := imports.Process(path, file.content, nil)
		if err != nil {
			return fmt.Errorf("formatGo file %s: %v", path, err)
		}
		if err := ioutil.WriteFile(path, src, 0644); err != nil {
			return fmt.Errorf("write file %s: %v", path, err)
		}
	}
	return nil
}

// formatDart runs "dartfmt" on all assets.
func (a assets) formatDart() error {
	args := []string{"-w"}
	for _, dir := range a.dirs {
		args = append(args, dir)
	}

	if err := exec.Command("dartfmt", args...).Run(); err != nil {
		if err == exec.ErrNotFound {
			// "dartfmt" is not available
			fmt.Println("The command 'dartfmt' was not found on your system. Generated code remains unformatted.")
		} else {
			return err
		}
	}

	return nil
}

// Dart type for a given go type.
func dartType(typeMappings []*TypeMapping) func(*field.TypeInfo) string {
	mappings := dartTypeNames

	for _, m := range typeMappings {
		mappings[m.Go] = m.Dart
	}

	return func(t *field.TypeInfo) string {
		if s, ok := mappings[t.String()]; ok {
			return s
		}

		return dartTypeNames["invalid"]
	}
}

// Extract the dart fields of a given type.
func dartRequestFields(dt func(*field.TypeInfo) string) func(*gen.Type, string) []dartField {
	return func(t *gen.Type, a string) []dartField {
		s := make([]dartField, 0)

		for _, f := range t.Fields {
			if f.Annotations["FieldGen"] == nil || (a != "" && !f.Annotations["FieldGen"].(map[string]interface{})[a].(bool)) {
				df := dartField{Type: dt(f.Type), Field: f}

				if f.HasGoType() {
					df.Converter = fmt.Sprintf("@%sConverter()", dt(f.Type))
				}

				s = append(s, df)
			}
		}

		for _, e := range t.Edges {
			skip := e.Type.Annotations["HandlerGen"] != nil && e.Type.Annotations["HandlerGen"].(map[string]interface{})["Skip"].(bool)
			include := e.Annotations["FieldGen"] == nil || (a != "" && !e.Annotations["FieldGen"].(map[string]interface{})[a].(bool))
			if !skip && include {
				t := dt(e.Type.ID.Type)
				if !e.Unique {
					t = fmt.Sprintf("List<%s>", t)
				}

				s = append(s, dartField{Type: t, Edge: e})
			}
		}

		return s
	}
}

// What edges to eager-load.
func eagerLoadedEdges(n *gen.Type, groupKey string) []*gen.Edge {
	r := make([]*gen.Edge, 0)

	if n.Annotations["HandlerGen"] != nil {
		if as, ok := n.Annotations["HandlerGen"].(map[string]interface{}); ok {
			if ls, ok := as[groupKey].([]interface{}); ok {
				for _, e := range n.Edges {
					if t, ok := reflect.StructTag(e.StructTag).Lookup("groups"); ok {
						gs := strings.Split(t, ",")
						for _, g := range ls {
							for _, g1 := range gs {
								if g == g1 {
									r = append(r, e)
								}
							}
						}
					}
				}
			}
		}
	}

	return r
}

func pkgImports(g *gen.Graph) []string {
	i := make(map[string]struct{})

	for _, n := range g.Nodes {
		for _, f := range n.Fields {
			if f.HasGoType() {
				i[f.Type.PkgPath] = struct{}{}
			}
		}
	}

	r := make([]string, 0)
	for k := range i {
		r = append(r, k)
	}

	return r
}

func dec(i int) int {
	return i - 1
}
