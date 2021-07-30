package elk

import (
	"entgo.io/ent/entc/gen"
	"fmt"
	"strings"
)

type (
	// EdgeToLoad specifies and edge to load for a type.
	EdgeToLoad struct {
		Edge        *gen.Edge
		EdgesToLoad []EdgeToLoad
		Groups      []string // TODO: Check if this isn't redundant ...
	}
	// EdgesToLoad is a list of several EdgeToLoad.
	EdgesToLoad []EdgeToLoad
)

// EntQuery simply runs EntQuery on every item in the list.
func (es EdgesToLoad) EntQuery() string {
	b := new(strings.Builder)
	for _, e := range es {
		b.WriteString(e.EntQuery())
	}
	return b.String()
}

// EntQuery constructs the code to eager load all the defined edges for the given edge.
func (etl EdgeToLoad) EntQuery() string {
	b := new(strings.Builder)

	b.WriteString(fmt.Sprintf(".%s(", strings.Title(etl.Edge.EagerLoadField())))
	for _, e := range etl.EdgesToLoad {
		b.WriteString(fmt.Sprintf("func (q *ent.%s) {\nq%s\n}", e.Edge.Owner.QueryName(), e.EntQuery()))
	}
	b.WriteString(")")

	return b.String()
}

func edgesToLoad(n *gen.Type, action string) (EdgesToLoad, error) {
	// If there are no annotations given do not load any edges.
	a := &SchemaAnnotation{}
	if n.Annotations == nil || n.Annotations[a.Name()] == nil {
		return nil, nil
	}

	// Decode the types annotation and extract the groups requested for the given action.
	if err := a.Decode(n.Annotations[a.Name()]); err != nil {
		return nil, err
	}

	var g []string
	switch action {
	case actionCreate:
		g = a.CreateGroups
	case actionRead:
		g = a.ReadGroups
	case actionUpdate:
		g = a.UpdateGroups
	case actionList:
		g = a.ListGroups
	}

	return edgesToLoadHelper(n, make(map[string]uint), g)
}

// edgesToLoadHelper recursively collects the edges to load on this type for requested groups on the given action.
func edgesToLoadHelper(n *gen.Type, visited map[string]uint, groupsToLoad []string) (EdgesToLoad, error) {
	// What edges to load on this type.
	edges := make(EdgesToLoad, 0)

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
			etl, err := edgesToLoadHelper(e.Type, visited, groupsToLoad)
			if err != nil {
				return nil, err
			}

			edges = append(edges, EdgeToLoad{
				Edge:        e,
				EdgesToLoad: etl,
				Groups:      groupsToLoad,
			})
		}
	}

	return edges, nil
}

func encodeTypeAndEdge(n *gen.Type, e *gen.Edge) string {
	return n.Name + "_" + e.Name
}
