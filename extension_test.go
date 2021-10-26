package elk

// import (
// 	"github.com/stretchr/testify/require"
// 	"testing"
// )
//
// func TestExtensionOption(t *testing.T) {
// 	ex, err := NewExtension()
// 	require.EqualError(t, err,
// 		`no generator enabled: enable one by providing either "GenerateSpec()" or "GenerateHandlers()" to "NewExtension()"`,
// 	)
// 	require.Nil(t, ex)
//
// 	ex, err = NewExtension(GenerateHandlers())
// 	require.NoError(t, err)
// 	require.Len(t, ex.hooks, 1)
//
// 	ex, err = NewExtension(GenerateSpec(""))
// 	require.EqualError(t, err, "spec filename cannot be empty")
// 	require.Nil(t, ex)
//
// 	ex, err = NewExtension(GenerateSpec("spec.json", SpecTitle("")))
// 	require.NoError(t, err)
// 	require.Len(t, ex.hooks, 1)
// 	require.Len(t, ex.specHooks, 1)
//
// 	ex, err = NewExtension(GenerateHandlers(), GenerateSpec("spec.json"))
// 	require.NoError(t, err)
// 	require.Len(t, ex.hooks, 2)
// }
