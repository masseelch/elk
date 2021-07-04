package elk

import (
	"entgo.io/ent/entc/gen"
	"fmt"
	"strings"
)

type (
	// eagerLoadEdges holds information about the edges to eager load on in ent query
	eagerLoadEdges struct {
		// what edges to eager load on the current query.
		edges []eagerLoadEdge
		// queryName is the type of the query of this edge.
		queryName string
	}
	// eagerLoadEdge holds information about a edge to eager load in an ent query.
	eagerLoadEdge struct {
		// method is the name of the method on the query builder to load this edge.
		method string
		// eagerLoadEdges enables recursive eager loading.
		eagerLoadEdges *eagerLoadEdges
	}
)

func (el eagerLoadEdges) Code(name ...string) string {
	n := "q"
	if len(name) > 0 {
		n = name[0]
	}

	b := new(strings.Builder)
	b.WriteString(n)

	for _, e := range el.edges {
		b.WriteString(fmt.Sprintf(".%s(", e.method))

		if e.eagerLoadEdges != nil {
			n += "_"
			b.WriteString(fmt.Sprintf("func (%s *ent.%s) {\n%s\n}", n, e.eagerLoadEdges.queryName, e.eagerLoadEdges.Code(n)))
		}

		b.WriteString(")")
	}

	return b.String()
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

	return edgesToLoadHelper(n, make(map[string]uint), g)
}

// edgesToLoadHelper returns the query to use to eager load all edges as required by the annotations defined
// on the given node.
func edgesToLoadHelper(n *gen.Type, visited map[string]uint, groupsToLoad []string) (*eagerLoadEdges, error) {
	// What edges to load on this type.
	edges := make([]eagerLoadEdge, 0)

	// Iterate over the edges of the given type.
	// If the type has an edge we need to eager load, do so.
	// Recursively go down the current edges edges and eager load those too.
	for _, e := range n.Edges {
		// Parse the edges annotation.
		a := Annotation{}
		if e.Annotations != nil && e.Annotations[a.Name()] != nil {
			if err := a.Decode(e.Annotations[a.Name()]); err != nil {
				return nil, err
			}
		}
		a.EnsureDefaults()

		// If we have reached the max depth on this field for the given type stop the recursion.
		if visited[encodeTypeAndEdge(n, e)] >= a.MaxDepth {
			continue
		}

		// This edge on this type has been visited. Increase counter.
		visited[encodeTypeAndEdge(n, e)]++

		// TODO: Take the DefaultOrder-Annotation into account.

		// If the edge has at least one of the groups requested load the edge.
		if a.Groups.Match(groupsToLoad) {
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

func encodeTypeAndEdge(n *gen.Type, e *gen.Edge) string {
	return n.Name + "_" + e.Name
}
