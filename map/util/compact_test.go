package maputil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmptyMap(t *testing.T) {
	res, err := Compact(map[string]interface{}{})
	require.Nil(t, err)
	require.Empty(t, res)
}

func TestNilMap(t *testing.T) {
	res, err := Compact(nil)
	require.Nil(t, err)
	require.Empty(t, res)
}

func TestMapWithNils(t *testing.T) {
	res, err := Compact(map[string]interface{}{
		"some": "string",
		"key":  nil,
	})
	require.Nil(t, err)
	require.Equal(t, map[string]interface{}{"some": "string"}, res)
}
