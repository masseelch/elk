package elk

import (
	"entgo.io/ent/entc/gen"
	"errors"
	"fmt"
	"strings"
)

const maxDepth = 25

type (
	// Edge specifies and edge to load for a type.
	Edge struct {
		*gen.Edge
		Edges Edges
	}
	// Edges is a list of multiple EdgeToLoad.
	Edges []Edge
	// walk is a node sequence in the schema graph. Used to keep track when computing EdgesToLoad.
	walk []string
)

// EntQuery simply runs EntQuery on every item in the list.
func (es Edges) EntQuery() string {
	b := new(strings.Builder)
	for _, e := range es {
		b.WriteString(e.EntQuery())
	}
	return b.String()
}

// EntQuery constructs the code to eager load all the defined edges for the given edge.
func (e Edge) EntQuery() string {
	b := new(strings.Builder)
	b.WriteString(fmt.Sprintf(".%s(", strings.Title(e.EagerLoadField())))
	if len(e.Edges) > 0 {
		b.WriteString(fmt.Sprintf("func (q *ent.%s) {\nq%s\n}", e.Type.QueryName(), e.Edges.EntQuery()))
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

// reachedMaxDepth returns if the walk has reached a depth greater than maxDepth.
func (w walk) reachedMaxDepth() bool {
	return len(w) > maxDepth
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

// edges returns the EdgesToLoad for the given node and operation.
func edges(n *gen.Type, a string) (Edges, error) {
	g, err := groupsForOperation(n, a)
	if err != nil {
		return nil, err
	}
	return edgesHelper(n, walk{}, g)
}

// edgesHelper recursively collects the edges to load on this type for requested groups on the given operation.
func edgesHelper(n *gen.Type, w walk, groupsToLoad []string) (Edges, error) {
	// If we have reached maxDepth there most possibly is an unwanted circular reference.
	if w.reachedMaxDepth() {
		return nil, errors.New(fmt.Sprintf("max depth of %d reached: ", maxDepth))
	}
	// What edges to load on this type.
	var es Edges
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
			es1, err := edgesHelper(e.Type, w, groupsToLoad)
			if err != nil {
				return nil, err
			}
			// Done visiting this node. Remove this node from our walk.
			w.pop()
			es = append(es, Edge{Edge: e, Edges: es1})
		}
	}
	return es, nil
}

func encodeTypeAndEdge(n *gen.Type, e *gen.Edge) string {
	return n.Name + "." + e.Name
}
