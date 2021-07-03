package elk

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEagerLoadEdges_Code(t *testing.T) {
	es := eagerLoadEdges{}
	require.Equal(t, "q", es.Code())

	es = eagerLoadEdges{
		edges: []eagerLoadEdge{
			{method: "WithEdgeOne"},
			{method: "WithEdgeTwo"},
		},
	}
	require.Equal(t, "a.WithEdgeOne().WithEdgeTwo()", es.Code("a"))

	es = eagerLoadEdges{
		edges: []eagerLoadEdge{
			{
				method: "WithEdgeOne",
				eagerLoadEdges: &eagerLoadEdges{
					queryName: "EdgeOneQuery",
					edges: []eagerLoadEdge{
						{method: "WithEdgeOneEdgeOne"},
						{method: "WithEdgeOneEdgeTwo"},
						{
							method: "WithEdgeOneEdgeThree",
							eagerLoadEdges: &eagerLoadEdges{
								queryName: "EdgeOneEdgeThreeQuery",
								edges:     []eagerLoadEdge{{method: "WithEdgeOneEdgeThreeEdgeOne"}},
							},
						},
					},
				},
			},
			{
				method: "WithEdgeTwo",
			},
		},
	}
	require.Equal(t, "q.WithEdgeOne(func (q_ *ent.EdgeOneQuery) {\nq_.WithEdgeOneEdgeOne().WithEdgeOneEdgeTwo().WithEdgeOneEdgeThree(func (q__ *ent.EdgeOneEdgeThreeQuery) {\nq__.WithEdgeOneEdgeThreeEdgeOne()\n})\n}).WithEdgeTwo()", es.Code())
}
