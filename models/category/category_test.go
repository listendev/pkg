package category

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCase(t *testing.T) {
	for _, c := range all {
		x := c.Case()
		assert.Equal(t, c.String(), x.Original())
	}
}

func TestFromUint64(t *testing.T) {
	_, err := FromUint64(math.MaxUint64)
	assert.Error(t, err)

	for _, c := range all {
		r, e := FromUint64(uint64(c))
		assert.Nil(t, e)
		assert.Equal(t, c, r)
	}
}

func TestFromString(t *testing.T) {
	for _, c := range all {
		a, err := FromString(c.String())
		assert.Nil(t, err)
		assert.Equal(t, c, a)

		b, err := FromString(string(c.Case()))
		assert.Nil(t, err)
		assert.Equal(t, b, a)
	}
}

func TestMarshal(t *testing.T) {
	for _, c := range all {
		r, e := json.Marshal(c)
		assert.Nil(t, e)
		assert.NotNil(t, r)
		assert.Equal(t, fmt.Sprintf("%q", c.Case()), string(r))
	}
}

func TestUnmarshal(t *testing.T) {
	for _, c := range all {
		var v1 Category
		e1 := json.Unmarshal([]byte(fmt.Sprintf("%q", c.String())), &v1)
		assert.Nil(t, e1)
		assert.Equal(t, c, v1)

		var v2 Category
		e2 := json.Unmarshal([]byte(fmt.Sprintf("%q", c.Case())), &v2)
		assert.Nil(t, e2)
		assert.Equal(t, c, v2)
	}
}
