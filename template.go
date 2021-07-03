package elk

import (
	"entgo.io/ent/entc/gen"
	"github.com/masseelch/elk/internal"
	"strings"
	"text/template"
)

//go:generate go run github.com/go-bindata/go-bindata/go-bindata -o=internal/bindata.go -pkg=internal -modtime=1 ./template/...

const (
	actionCreate = "create"
	actionRead   = "read"
	actionUpdate = "update"
	actionDelete = "delete"
)

var (
	// HTTPTemplates holds all templates for generating http handlers.
	HTTPTemplates = []*gen.Template{
		parse("template/http/handler.tmpl"),
		parse("template/http/create.tmpl"),
		parse("template/http/read.tmpl"),
	}

	// TemplateFuncs contains the extra template functions used by elk.
	TemplateFuncs = template.FuncMap{
		"edgesToLoad": edgesToLoad,
	}
)

func parse(path string) *gen.Template {
	return gen.MustParse(gen.NewTemplate(path).
		Funcs(gen.Funcs).
		Funcs(TemplateFuncs).
		Parse(string(internal.MustAsset(path))))
}

// edgesToLoad generates the code to eager load as defined by the elk annotation.
func edgesToLoad(n *gen.Type, action string) (*eagerLoadEdges, error) {
	// If there are no annotations given do not load any edges.
	a := &SchemaAnnotation{}
	if n.Annotations == nil || n.Annotations[a.Name()] == nil {
		return nil, nil
	}

	// Load the annotation.
	if err := a.Decode(n.Annotations[a.Name()]); err != nil {
		return nil, err
	}

	// Extract the groups requested.
	var g []string
	switch action {
	case actionCreate:
		g = a.CreateGroups
	case actionRead:
		g = a.ReadGroups
	case actionUpdate:
		g = a.UpdateGroups
	case actionDelete:
		g = a.DeleteGroups
	}

	return edgesToLoadHelper(n, make(map[string]bool), g)
}

// edgesToLoadHelper returns the query to use to eager load all edges as required by the annotations defined
// on the given node.
func edgesToLoadHelper(n *gen.Type, visited map[string]bool, groupsToLoad []string) (*eagerLoadEdges, error) {
	// Stop recursion (termination condition)
	if groupsToLoad == nil || visited[n.Name] {
		return nil, nil
	}

	// Mark this node as visited.
	visited[n.Name] = true

	// What edges to load on this type.
	edges := make([]eagerLoadEdge, 0)

	// Iterate over the edges of the given type.
	// If the type has an edge we need to eager load, do so.
	// Recursively go down the current edges edges and eager load those too.
	for _, e := range n.Edges {
		// Do not include a type we already saw above this branch of the serialization tree.
		if visited[e.Type.Name] {
			continue
		}

		// TODO: Take the DefaultOrder-Annotation into account.

		// Groups defined on the current edge.
		gs := Groups{}
		a := Annotation{}
		if e.Annotations != nil && e.Annotations[a.Name()] != nil {
			if err := a.Decode(e.Annotations[a.Name()]); err != nil {
				return nil, err
			}

			gs = a.Groups
		}

		// If the edge has at least one of the groups requested load the edge.
		if gs.Match(groupsToLoad) {
			// Recursively collect the eager loads of this edges edges.
			eagerLoadEdges, err := edgesToLoadHelper(e.Type, visited, groupsToLoad)
			if err != nil {
				return nil, err
			}

			edges = append(edges, eagerLoadEdge{
				eagerLoadEdges: eagerLoadEdges,
				method:         strings.Title(e.EagerLoadField()),
			})
		}
	}

	// If there are no edges to load on this type return nil.
	if len(edges) == 0 {
		return nil, nil
	}

	return &eagerLoadEdges{
		edges:     edges,
		queryName: n.QueryName(),
	}, nil
}
