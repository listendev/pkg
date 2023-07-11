package ecosystem

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromUint64(t *testing.T) {
	_, err := FromUint64(math.MaxUint64)
	assert.Error(t, err)

	for _, c := range all {
		r, e := FromUint64(uint64(c))
		assert.Nil(t, e)
		assert.Equal(t, c, r)
	}

	_, e := FromUint64(0)
	assert.Error(t, e)
}

func TestFromString(t *testing.T) {
	for _, c := range all {
		a, err := FromString(c.String())
		assert.Nil(t, err)
		assert.Equal(t, c, a)

		b, err := FromString(c.Case())
		assert.Nil(t, err)
		assert.Equal(t, b, a)

		upper, err := FromString(strings.ToUpper(c.String()))
		assert.Nil(t, err)
		assert.Equal(t, b, upper)
	}

	_, e := FromString("none")
	assert.Error(t, e)
}

func TestMarshal(t *testing.T) {
	for _, c := range all {
		r, e := json.Marshal(c)
		assert.Nil(t, e)
		assert.NotNil(t, r)
		assert.Equal(t, fmt.Sprintf("%q", c.Case()), string(r))
	}

	_, err1 := json.Marshal(Ecosystem(math.MaxUint64))
	assert.Error(t, err1)

	_, err2 := json.Marshal(None)
	assert.Error(t, err2)
}

func TestUnmarshal(t *testing.T) {
	for _, x := range all {
		var v1 Ecosystem
		e1 := json.Unmarshal([]byte(fmt.Sprintf("%q", x.String())), &v1)
		if assert.Nil(t, e1) {
			assert.Equal(t, x, v1)
		}

		var v2 Ecosystem
		e2 := json.Unmarshal([]byte(fmt.Sprintf("%q", x.Case())), &v2)
		if assert.Nil(t, e2) {
			assert.Equal(t, x, v2)
		}
	}

	var e Ecosystem
	assert.Error(t, json.Unmarshal([]byte(fmt.Sprintf("%q", "none")), &e))
}
