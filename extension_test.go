package elk

import (
	"github.com/masseelch/elk/openapi"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExtensionOption(t *testing.T) {
	c := EasyJsonConfig{
		NoStdMarshalers:          false,
		SnakeCase:                true,
		LowerCamelCase:           false,
		OmitEmpty:                true,
		DisallowUnknownFields:    false,
		SkipMemberNameUnescaping: true,
	}
	spec, err := openapi.New()
	ex, err := NewExtension(WithEasyJsonConfig(c), WithOpenAPISpec(spec))
	require.NoError(t, err)
	require.Equal(t, c, ex.easyjsonConfig)
	require.Equal(t, spec, ex.openAPISpec)
}
