package elk

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

	require.Equal(t, `groups:"group_one,GROUP_two,group:3"`, Groups{"group_one", "GROUP_two", "group:3"}.StructTag())
}
