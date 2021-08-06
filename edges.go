package elk

import (
	"entgo.io/ent/entc/gen"
	"errors"
	"fmt"
	"strings"
)

const maxDepth = 25

type (
	// EdgeToLoad specifies and edge to load for a type.
	EdgeToLoad struct {
		Edge        *gen.Edge
		EdgesToLoad EdgesToLoad
		Groups      []string
	}
	// EdgesToLoad is a list of several EdgeToLoad.
	EdgesToLoad []EdgeToLoad
	// walk is a node sequence in the schema graph. Used to keep track when computing EdgesToLoad.
	walk []string
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

	if len(etl.EdgesToLoad) > 0 {
		b.WriteString(fmt.Sprintf("func (q *ent.%s) {\nq%s\n}", etl.Edge.Type.QueryName(), etl.EdgesToLoad.EntQuery()))
	}

	b.WriteString(")")

	return b.String()
}

// cycleDepth determines the length of a cycle on the last visited node.
//   <nil>: 0 -> no visits at all
// a->b->c: 1 -> 1st visit on c
// a->b->b: 2 -> 2nd visit on b
// a->a->a: 3 -> 3rd visit on a
// a->b->a: 2 -> 2nd visit on a
func (w walk) cycleDepth() uint {
	if len(w) == 0 {
		return 0
	}
	n := w[len(w)-1]
	c := uint(1)
	for i := len(w) - 2; i >= 0; i-- {
		if n == w[i] {
			c++
		}
	}
	return c
}

// reachedMaxDepth returns if the walk has reached a depth greater then maxDepth.
func (w walk) reachedMaxDepth() bool {
	return len(w) > maxDepth
}

// tail returns a formatted string of the last c steps in walk.
func (w walk) tail(c int) string {
	f := "%s"
	for i := 0; i < c; i++ {
		f += " -> %s"
	}
	if c > len(w) {
		c = len(w)
	}
	return fmt.Sprintf(f, w[len(w)-c:])
}

// push adds a new step to the walk.
func (w *walk) push(s string) {
	*w = append(*w, s)
}

// pop removed the last step of the walk.
func (w *walk) pop() {
	if len(*w) > 0 {
		*w = (*w)[:len(*w)-1]
	}
}

// edgesToLoad returns the EdgesToLoad for the given node and action.
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

	return edgesToLoadHelper(n, walk{}, g)
}

// edgesToLoadHelper recursively collects the edges to load on this type for requested groups on the given action.
func edgesToLoadHelper(n *gen.Type, w walk, groupsToLoad []string) (EdgesToLoad, error) {
	// If we have reached maxDepth there most possibly is an unwanted circular reference.
	if w.reachedMaxDepth() {
		return nil, errors.New(fmt.Sprintf("max depth of %d reached: ", maxDepth))
	}

	// What edges to load on this type.
	var edges []EdgeToLoad

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

		// If the edge has at least one of the groups requested load the edge.
		if a.Groups.Match(groupsToLoad) {
			// Add the current step to our walk, since we will add this edge.
			w.push(encodeTypeAndEdge(n, e))

			// If we have reached the max depth on this field for the given type stop the recursion. Backtrack!
			if w.cycleDepth() > a.MaxDepth {
				w.pop()
				continue
			}

			// Recursively collect the eager loads of this edges edges.
			etl, err := edgesToLoadHelper(e.Type, w, groupsToLoad)
			if err != nil {
				return nil, err
			}

			// Done visiting this node. Remove this node from our walk.
			w.pop()

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
	return n.Name + "." + e.Name
}
