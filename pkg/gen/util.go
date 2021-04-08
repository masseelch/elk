package gen

import (
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"fmt"
	"golang.org/x/tools/imports"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"unicode"
)

var (
	dartTypeNames = map[string]string{
		"invalid":   "dynamic",
		"bool":      "bool",
		"time.Time": "DateTime",
		// "JSON":    "Map<String, dynamic>",
		// "UUID":    "String",
		// "bytes":   "dynamic",
		"enum":     "String",
		"string":   "String",
		"int":      "int",
		"int8":     "int",
		"int16":    "int",
		"int32":    "int",
		"int64":    "int",
		"uint":     "int",
		"uint8":    "int",
		"uint16":   "int",
		"uint32":   "int",
		"uint64":   "int",
		"float32":  "double",
		"float64":  "double",
		"[]string": "List<String>",
	}
)

type (
	dartFields []dartField
	dartField  struct {
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
		Source    string
		Target    string
		Package   string
		Templates []string
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

func (d dartField) Immutable() bool {
	if d.IsEdge() {
		return false
	}

	return d.Field.Immutable
}

func (d dartField) StructField() string {
	if d.IsEdge() {
		return d.Edge.StructField()
	}

	return d.Field.StructField()
}

func (d dartFields) String() string {
	b := new(strings.Builder)

	b.WriteString("[")

	for i, df := range d {
		if i != 0 {
			b.WriteString("; ")
		}

		if df.IsEdge() {
			b.WriteString(fmt.Sprintf("edge: %s, type: %s, conv: %s", df.Edge.Name, df.Type, df.Converter))
		} else {
			b.WriteString(fmt.Sprintf("field: %s, type: %s, conv: %s", df.Field.Name, df.Type, df.Converter))
		}
	}

	return b.String()
}

func (d dartFields) ConverterFor(f *gen.Field) string {
	if f == nil {
		return ""
	}

	for _, df := range d {
		if df.IsEdge() {
			continue
		}

		if df.Field.Name == f.Name {
			return df.Converter
		}
	}

	return ""
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
func (a assets) formatDart(dirs []string) error {
	args := []string{"-w", "--line-length=120"}
	for _, dir := range dirs {
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

		if t.Type == field.TypeJSON {
			fmt.Println(t)
		}

		// Try to guess the type. dynamic otherwise.
		return dartTypeNames[t.Type.String()]
	}
}

// Extract the dart fields of a given type.
func dartRequestFields(c *FlutterConfig, dt func(*field.TypeInfo) string) func(*gen.Type, string) dartFields {
	return func(t *gen.Type, a string) dartFields {
		s := make(dartFields, 0)

		for _, f := range t.Fields {
			if f.Annotations["FieldGen"] == nil || a == "" || (a != "" && !f.Annotations["FieldGen"].(map[string]interface{})[a].(bool)) {
				df := dartField{Type: dt(f.Type) + "?", Field: f}

				if f.Annotations["FieldGen"] != nil && f.Annotations["FieldGen"].(map[string]interface{})["MapGoType"].(bool) && f.HasGoType() {
					// Find the Type-Mapping. If a converter is needed use it.
					for _, tm := range c.TypeMappings {
						if tm.Go == f.Type.String() && tm.ConverterNeeded {
							df.Converter = fmt.Sprintf("@%sConverter()", dt(f.Type))
						}
					}
				}

				s = append(s, df)
			}
		}

		for _, e := range t.Edges {
			skip := e.Type.Annotations["HandlerGen"] != nil && e.Type.Annotations["HandlerGen"].(map[string]interface{})["Skip"].(bool)
			include := e.Annotations["FieldGen"] == nil || a == "" || (a != "" && !e.Annotations["FieldGen"].(map[string]interface{})[a].(bool))
			if !skip && include {
				t := dt(e.Type.ID.Type)
				if !e.Unique {
					t = fmt.Sprintf("List<%s?>", t)
				}

				s = append(s, dartField{Type: t + "?", Edge: e})
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

// What edges to eager-load.
func eagerLoadQueryByGroupKeyRecursive(n *gen.Type, groupKey string, builderVar string, visitedTypes []*gen.Type, groups []string) *string {
	// Use the given groups. If there are non given this is the first leaf in the serialization tree.
	// In that case use the types annotation to fetch the groups.

	annotationGroups := groups
	if annotationGroups == nil || len(annotationGroups) == 0 {
		annotationGroups = extractGroupsFromAnnotations(n, groupKey)
	}

	// Stop the recursion
	if annotationGroups == nil || typeVisited(visitedTypes, n) {
		return nil
	}

	// Mark this node as visited.
	visitedTypes = append(visitedTypes, n)

	// Start the query
	b := new(strings.Builder)
	b.WriteString(builderVar)

	// Iterate over the edges of the given type.
	// If the type has an edge we need to eager load do so.
	// Recursively go down the current edges edges and eager load those too.
	for _, e := range n.Edges {
		// Do not include a type we already saw above this branch of the serialization tree.
		if typeVisited(visitedTypes, e.Type) {
			continue
		}

		// TODO: Take the DefaultOrder-Annotation into account.

		// Get the groups of the edge out of its group-tag.
		tagGroups := extractGroupsFromTag(e.StructTag)

		// If the annotationGroups contain a group that is present on the edge eager load the edge.
		if haveSameGroup(annotationGroups, tagGroups) {
			b.WriteString(fmt.Sprintf(".With%s(", gen.Funcs["pascal"].(func(string) string)(e.Name)))

			// Recursively collect the eager loads of this edges edges.
			subQueryBuilderVar := "_" + builderVar
			if q := eagerLoadQueryByGroupKeyRecursive(e.Type, groupKey, "_"+builderVar, visitedTypes, annotationGroups); q != nil {
				b.WriteString(fmt.Sprintf("func(%s *ent.%s){\n%s\n}", subQueryBuilderVar, e.Type.QueryName(), *q))
			}

			// Finally close the
			b.WriteString(")")
		}
	}

	s := b.String()
	if s == builderVar {
		return nil
	}

	return &s
}

func extractGroupsFromAnnotations(n *gen.Type, groupKey string) []string {
	if n.Annotations["HandlerGen"] != nil {
		if annos, ok := n.Annotations["HandlerGen"].(map[string]interface{}); ok {
			if groupsI, ok := annos[groupKey].([]interface{}); ok {
				groups := make([]string, len(groupsI))
				for i := range groupsI {
					groups[i] = groupsI[i].(string)
				}

				return groups
			}
		}
	}

	return nil
}

func haveSameGroup(a []string, b []string) bool {
	for _, ae := range a {
		for _, be := range b {
			if ae == be {
				return true
			}
		}
	}

	return false
}

func extractGroupsFromTag(tag string) []string {
	if t, ok := reflect.StructTag(tag).Lookup("groups"); ok {
		return strings.Split(t, ",")
	}

	return nil
}

func typeVisited(visited []*gen.Type, t *gen.Type) bool {
	for _, _t := range visited {
		if _t == t {
			return true
		}
	}

	return false
}

func visited(v ...*gen.Type) []*gen.Type {
	return v
}

func pkgImports(g *gen.Graph) []string {
	i := make(map[string]struct{})

	for _, n := range g.Nodes {
		if n.ID.HasGoType() {
			i[n.ID.Type.PkgPath] = struct{}{}
		}
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

func lowerFirst(s string) string {
	return string(unicode.ToLower(rune(s[0]))) + s[1:]
}
