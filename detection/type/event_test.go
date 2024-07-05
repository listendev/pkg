package detectiontype

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventTypesAsStringsOptions(t *testing.T) {
	require.Contains(t, EventTypesAsStrings(), ShellConfigModification.String())
	require.Contains(t, EventTypesAsStrings(ApplyCase), SudoersModification.Case())
	require.Contains(t, EventTypesAsStrings(ApplyCase, SingleQuotes), fmt.Sprintf("'%s'", SslCertificateAccess.Case()))
	require.Contains(t, EventTypesAsStrings(SingleQuotes, WithValue), fmt.Sprintf("'%s' = %d", ProcessCodeModification, ProcessCodeModification))
}

func TestFromUint64(t *testing.T) {
	_, err := FromUint64(math.MaxUint64)
	assert.Error(t, err)

	for _, c := range events {
		r, e := FromUint64(uint64(c))
		assert.Nil(t, e)
		assert.Equal(t, c, r)
	}

	_, e := FromUint64(0)
	assert.Error(t, e)
}

func TestFromString(t *testing.T) {
	for _, e := range events {
		s := e.String()
		a, err := FromString(s)
		assert.Nil(t, err)
		assert.Equal(t, e, a)

		c := e.Case()
		b, err := FromString(c)
		assert.Nil(t, err)
		assert.Equal(t, b, a)

		upper, err := FromString(strings.ToUpper(s))
		assert.Nil(t, err)
		assert.Equal(t, b, upper)
	}

	_, e1 := FromString("none")
	require.Error(t, e1)
	_, e2 := FromString("None")
	require.Error(t, e2)
}
