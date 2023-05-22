package verdictcode

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/garnet-org/pkg/analysisrequest"
	"github.com/stretchr/testify/assert"
)

func TestFromUint64(t *testing.T) {
	c, err := FromUint64(math.MaxUint64, false)
	assert.Error(t, err)
	assert.Equal(t, UNK, c)

	fin001, err := FromUint64(uint64(FNI001), false)
	assert.Nil(t, err)
	assert.Equal(t, FNI001, fin001)

	ddn01, err := FromUint64(uint64(DDN01), false)
	assert.Nil(t, err)
	assert.Equal(t, DDN01, ddn01)

	tsn01, err := FromUint64(uint64(TSN01), false)
	assert.Nil(t, err)
	assert.Equal(t, TSN01, tsn01)
}

func TestFromString(t *testing.T) {
	c, err := FromString("CIAO", false)
	assert.Error(t, err)
	assert.Equal(t, UNK, c)

	fin001, err := FromString("FNI001", false)
	assert.Nil(t, err)
	assert.Equal(t, FNI001, fin001)

	ddn01, err := FromString("DDN01", false)
	assert.Nil(t, err)
	assert.Equal(t, DDN01, ddn01)

	tsn01, err := FromString("TSN01", false)
	assert.Nil(t, err)
	assert.Equal(t, TSN01, tsn01)
}

func TestGetBy(t *testing.T) {
	_, err := GetBy(analysisrequest.Nop)
	assert.NotNil(t, err)

	codes, err := GetBy(analysisrequest.NPMDepsDev)
	assert.Nil(t, err)
	assert.NotNil(t, codes)
}

func TestMarshal(t *testing.T) {
	r1, e1 := json.Marshal(FNI001)
	assert.Nil(t, e1)
	assert.NotNil(t, r1)
	assert.Equal(t, fmt.Sprintf("%q", FNI001.String()), string(r1))
}

func TestUnmarshal(t *testing.T) {
	var v1 Code
	e1 := json.Unmarshal([]byte(`"FNI001"`), &v1)
	assert.Nil(t, e1)
	assert.Equal(t, FNI001, v1)

	var v2 Code
	e2 := json.Unmarshal([]byte(`"UNK"`), &v2)
	if assert.NotNil(t, e2) {
		assert.Error(t, e2)
	}
}
