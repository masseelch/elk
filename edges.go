package elk

import (
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
