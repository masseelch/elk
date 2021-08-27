package serialization

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGroups(t *testing.T) {
	gs := Groups{}

	gs.Add("group")
	require.Len(t, gs, 1)
	require.True(t, gs.HasGroup("group"))
	require.False(t, gs.HasGroup("none"))

	gs.Add("group_1", "group_2")
	require.Len(t, gs, 3)
	require.True(t, gs.HasGroup("group"))
	require.True(t, gs.HasGroup("group_1"))
	require.True(t, gs.HasGroup("group_2"))
	require.False(t, gs.HasGroup("none"))

	require.False(t, gs.Match(Groups{"none", "nobody"}))
	require.True(t, gs.Match(Groups{"group", "nobody"}))

	require.True(t, Groups{}.Equal(Groups{}))
	require.True(t, Groups{"a"}.Equal(Groups{"a"}))
	require.True(t, Groups{"a", "b"}.Equal(Groups{"a", "b"}))
	require.False(t, Groups{"a"}.Equal(Groups{}))
	require.False(t, Groups{"a"}.Equal(Groups{"b"}))

	require.Equal(t, Groups{}.Hash(), Groups{}.Hash())
	require.Equal(t, Groups{"a"}.Hash(), Groups{"a"}.Hash())
	require.Equal(t, Groups{"a", "b"}.Hash(), Groups{"a", "b"}.Hash())
	require.NotEqual(t, Groups{"a", "b"}.Hash(), Groups{"ab"}.Hash())
	require.NotEqual(t, Groups{"a"}.Hash(), Groups{"b"}.Hash())

	c := Collection{{"a"}, {"a", "b"}}
	require.True(t, c.Contains(Groups{"a"}))
	require.True(t, c.Contains(Groups{"a", "b"}))
	require.False(t, c.Contains(Groups{"b"}))
}
