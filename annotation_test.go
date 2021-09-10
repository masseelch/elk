package elk

import (
	"github.com/masseelch/elk/serialization"
	"github.com/masseelch/elk/spec"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAnnotation(t *testing.T) {
	t.Parallel()

	a := CreateGroups("create", "groups")
	require.Equal(t, a.CreateGroups, serialization.Groups{"create", "groups"})

	a = ReadGroups("read", "groups")
	require.Equal(t, a.ReadGroups, serialization.Groups{"read", "groups"})

	a = UpdateGroups("update", "groups")
	require.Equal(t, a.UpdateGroups, serialization.Groups{"update", "groups"})

	a = ListGroups("list", "groups")
	require.Equal(t, a.ListGroups, serialization.Groups{"list", "groups"})

	a = CreatePolicy(Exclude)
	require.Equal(t, a.CreatePolicy, Exclude)

	a = ReadPolicy(Exclude)
	require.Equal(t, a.ReadPolicy, Exclude)

	a = UpdatePolicy(Exclude)
	require.Equal(t, a.UpdatePolicy, Exclude)

	a = DeletePolicy(Exclude)
	require.Equal(t, a.DeletePolicy, Exclude)

	a = ListPolicy(Exclude)
	require.Equal(t, a.ListPolicy, Exclude)

	a = SchemaPolicy(Expose)
	require.Equal(t, a.CreatePolicy, Expose)
	require.Equal(t, a.ReadPolicy, Expose)
	require.Equal(t, a.UpdatePolicy, Expose)
	require.Equal(t, a.DeletePolicy, Expose)
	require.Equal(t, a.ListPolicy, Expose)

	sec := spec.Security{{"apiKeySample": {}}}
	a = CreateSecurity(sec)
	require.Equal(t, a.CreateSecurity, sec)

	a = ReadSecurity(sec)
	require.Equal(t, a.ReadSecurity, sec)

	a = UpdateSecurity(sec)
	require.Equal(t, a.UpdateSecurity, sec)

	a = DeleteSecurity(sec)
	require.Equal(t, a.DeleteSecurity, sec)

	a = ListSecurity(sec)
	require.Equal(t, a.ListSecurity, sec)

	a = SchemaSecurity(sec)
	require.Equal(t, a.CreateSecurity, sec)
	require.Equal(t, a.ReadSecurity, sec)
	require.Equal(t, a.UpdateSecurity, sec)
	require.Equal(t, a.DeleteSecurity, sec)
	require.Equal(t, a.ListSecurity, sec)
}
