package elk

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestEdgesToLoad(t *testing.T) {
	// Load a graph.
	wd, err := os.Getwd()
	require.NoError(t, err)
	g, err := entc.LoadGraph(filepath.Join(wd, "internal", "petstore", "ent", "schema"), &gen.Config{
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
	q, err := edgesToLoad(p, actionRead)
	require.NoError(t, err)

	require.Equal(t, &eagerLoadEdges{
		queryName: "PetQuery",
		edges: []eagerLoadEdge{
			{
				method:         "WithOwner",
				eagerLoadEdges: nil,
			},
		},
	}, q)
}
