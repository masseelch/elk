package elk

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestWalk_CycleDepth(t *testing.T) {
	w := walk{}
	require.Equal(t, uint(0), w.cycleDepth())

	w = walk{"a"}
	require.Equal(t, uint(1), w.cycleDepth())

	w = walk{"a", "b"}
	require.Equal(t, uint(1), w.cycleDepth())

	w = walk{"a", "a"}
	require.Equal(t, uint(2), w.cycleDepth())

	w = walk{"a", "b", "b"}
	require.Equal(t, uint(2), w.cycleDepth())

	w = walk{"a", "b", "b", "c"}
	require.Equal(t, uint(1), w.cycleDepth())

	w = walk{"a", "b", "b", "a"}
	require.Equal(t, uint(2), w.cycleDepth())

	w = walk{"a", "a", "b", "a"}
	require.Equal(t, uint(3), w.cycleDepth())
}

func TestWalk_Push(t *testing.T) {
	w := walk{}
	require.Equal(t, walk{}, w)

	w.push("a")
	require.Equal(t, walk{"a"}, w)

	w.push("b")
	require.Equal(t, walk{"a", "b"}, w)
}

func TestWalk_Pop(t *testing.T) {
	w := walk{"a", "b", "c"}

	w.pop()
	require.Equal(t, walk{"a", "b"}, w)

	w.pop()
	require.Equal(t, walk{"a"}, w)

	w.pop()
	require.Equal(t, walk{}, w)
}

func TestEdges(t *testing.T) {
	// Load a graph.
	wd, err := os.Getwd()
	require.NoError(t, err)
	g, err := entc.LoadGraph(
		filepath.Join(wd, "internal", "simple", "ent", "schema"),
		&gen.Config{Templates: []*gen.Template{HTTPTemplate}},
	)
	require.NoError(t, err)

	// Generate the query to eager load edges for a read operation on a pet.
	var p *gen.Type
	for _, n := range g.Nodes {
		if n.Name == "Pet" {
			p = n
			break
		}
	}
	etls, err := edges(p, actionRead)
	require.NoError(t, err)

	require.Equal(
		t,
		".WithOwner()"+
			".WithFriends(func (q *ent.PetQuery) {\n"+
			"q.WithOwner().WithFriends(func (q *ent.PetQuery) {\n"+
			"q.WithOwner().WithFriends(func (q *ent.PetQuery) {\n"+
			"q.WithOwner()\n"+
			"})\n"+
			"})\n"+
			"})",
		etls.EntQuery(),
	)
}
