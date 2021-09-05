package elk

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExtensionOption(t *testing.T) {
	ex, err := NewExtension()
	require.EqualError(t, err,
		`no generator enabled: enable one by providing either "EnableSpecGenerator()" or "EnableHandlerGenerator()" to "NewExtension()"`,
	)
	require.Nil(t, ex)

	ex, err = NewExtension(EnableHandlerGenerator())
	require.NoError(t, err)
	require.Len(t, ex.hooks, 1)

	ex, err = NewExtension(EnableSpecGenerator(""))
	require.EqualError(t, err, "spec filename cannot be empty")
	require.Nil(t, ex)

	ex, err = NewExtension(EnableSpecGenerator("spec.json", SpecTitle("")))
	require.NoError(t, err)
	require.Len(t, ex.hooks, 1)
	require.Len(t, ex.specHooks, 1)

	ex, err = NewExtension(EnableHandlerGenerator(), EnableSpecGenerator("spec.json"))
	require.NoError(t, err)
	require.Len(t, ex.hooks, 2)
}
