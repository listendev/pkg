package detectiontype

import (
	"encoding/json"
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

func TestMarshal(t *testing.T) {
	for _, c := range events {
		r, e := json.Marshal(c)
		assert.Nil(t, e)
		assert.NotNil(t, r)
		assert.Equal(t, fmt.Sprintf("%q", c.Case()), string(r))
	}

	_, err1 := json.Marshal(Event(math.MaxUint64))
	assert.Error(t, err1)

	_, err2 := json.Marshal(None)
	assert.Error(t, err2)
}

func TestUnmarshal(t *testing.T) {
	for _, x := range events {
		var v1 Event
		e1 := json.Unmarshal([]byte(fmt.Sprintf("%q", x.String())), &v1)
		if assert.Nil(t, e1) {
			assert.Equal(t, x, v1)
		}

		var v2 Event
		e2 := json.Unmarshal([]byte(fmt.Sprintf("%q", x.Case())), &v2)
		if assert.Nil(t, e2) {
			assert.Equal(t, x, v2)
		}
	}

	var e Event
	require.Error(t, json.Unmarshal([]byte(fmt.Sprintf("%q", "none")), &e))
}

func TestUnmarshalText(t *testing.T) {
	for _, x := range events {
		var evt1 Event
		err1 := evt1.UnmarshalText([]byte(x.String()))
		require.Nil(t, err1)
		require.Equal(t, evt1, x)

		var evt2 Event
		err2 := evt2.UnmarshalText([]byte(x.Case()))
		require.Nil(t, err2)
		require.Equal(t, evt2, x)
	}
}
