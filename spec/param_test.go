package spec

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParameterPlace_MarshalJSON(t *testing.T) {
	for e, p := range map[string]ParameterPlace{
		"query":  query,
		"header": header,
		"path":   path,
		"cookie": cookie,
	} {
		j, err := json.Marshal(p)
		require.NoError(t, err)
		require.Equal(t, []byte(fmt.Sprintf(`"%s"`, e)), j)
	}
	_, err := json.Marshal(ParameterPlace(10))
	require.Error(t, err)
}
