package elk

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWithEasyJsonConfig(t *testing.T) {
	c := EasyJsonConfig{
		NoStdMarshalers:          false,
		SnakeCase:                true,
		LowerCamelCase:           false,
		OmitEmpty:                true,
		DisallowUnknownFields:    false,
		SkipMemberNameUnescaping: true,
	}
	ex, err := NewExtension(WithEasyJsonConfig(c))
	require.NoError(t, err)
	require.Equal(t, c, ex.easyjsonConfig)
}
