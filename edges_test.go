package elk

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestEagerLoadEdges_String(t *testing.T) {
	es := eagerLoadEdges{}
	require.Equal(t, "q", es.String())

	es = eagerLoadEdges{
		edges: []eagerLoadEdge{
			{method: "WithEdgeOne"},
			{method: "WithEdgeTwo"},
		},
	}
	require.Equal(t, "q.WithEdgeOne().WithEdgeTwo()", es.String())

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
	require.Equal(t, "q.WithEdgeOne(func (q_ *ent.EdgeOneQuery) {\nq_.WithEdgeOneEdgeOne().WithEdgeOneEdgeTwo().WithEdgeOneEdgeThree(func (q__ *ent.EdgeOneEdgeThreeQuery) {\nq__.WithEdgeOneEdgeThreeEdgeOne()\n})\n}).WithEdgeTwo()", es.String())
}

func TestEdgesToLoad(t *testing.T) {
	// Load a graph.
	wd, err := os.Getwd()
	require.NoError(t, err)
	g, err := entc.LoadGraph(filepath.Join(wd, "internal", "integration", "pets", "ent", "schema"), &gen.Config{
		Templates: HTTPTemplates,
		Hooks: []gen.Hook{
			AddGroupsTag,
		},
	})
	require.NoError(t, err)

	// Generate the query to eager load edges for a read operation on a pet.
	var p *gen.Type
	for _, n := range g.Nodes {
		if n.Name == "Pet" {
			p = n
			break
		}
	}
	es, err := edgesToLoad(p, actionRead)
	require.NoError(t, err)

	// Max-Depth of 3
	require.Equal(t, &eagerLoadEdges{
		queryName: "PetQuery",
		edges: []eagerLoadEdge{
			{
				method: "WithOwner",
			},
			{
				method: "WithFriends",
				eagerLoadEdges: &eagerLoadEdges{
					queryName: "PetQuery",
					edges: []eagerLoadEdge{{
						method: "WithFriends",
						eagerLoadEdges: &eagerLoadEdges{
							queryName: "PetQuery",
							edges: []eagerLoadEdge{{
								method: "WithFriends",
							}},
						},
					}},
				},
			},
		},
	}, es)
}
