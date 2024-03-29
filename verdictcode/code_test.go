package verdictcode

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/listendev/pkg/analysisrequest"
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
	assert.False(t, ddn01.UniquelyIdentifies())

	tsn01, err := FromUint64(uint64(TSN01), false)
	assert.Nil(t, err)
	assert.Equal(t, TSN01, tsn01)
	assert.True(t, tsn01.UniquelyIdentifies())
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

	codes, err := GetBy(analysisrequest.NPMAdvisory)
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

func TestRetrievingTheType(t *testing.T) {
	t1, e1 := FNI001.Type(false)
	assert.Nil(t, e1)
	assert.NotNil(t, t1)
	assert.Equal(t, analysisrequest.NPMInstallWhileDynamicInstrumentation, t1)

	t2, e2 := STN004.Type(false)
	assert.Nil(t, e2)
	assert.NotNil(t, t2)
	assert.Equal(t, analysisrequest.NPMStaticAnalysisEvalBase64, t2)

	t3, e3 := MDN03.Type(false)
	assert.Nil(t, e3)
	assert.NotNil(t, t3)
	assert.Equal(t, analysisrequest.NPMMetadataVersion, t3)
	assert.True(t, MDN03.UniquelyIdentifies())

	t4, e4 := MDN02.Type(false)
	assert.Nil(t, e4)
	assert.NotNil(t, t4)
	assert.Equal(t, analysisrequest.NPMMetadataVersion, t4)

	t5, e5 := MDN06.Type(false)
	assert.Nil(t, e5)
	assert.NotNil(t, t5)
	assert.Equal(t, analysisrequest.NPMMetadataMismatches, t5)

	t6, e6 := TSP01.Type(false)
	assert.Nil(t, e6)
	assert.NotNil(t, t6)
	assert.Equal(t, analysisrequest.PypiTyposquat, t6)

	t7, e7 := MDN04.Type(false)
	assert.Nil(t, e7)
	assert.NotNil(t, t7)
	assert.Equal(t, analysisrequest.NPMMetadataMaintainersEmailCheck, t7)

	t8, e8 := MDP04.Type(false)
	assert.Nil(t, e8)
	assert.NotNil(t, t8)
	assert.Equal(t, analysisrequest.PypiMetadataMaintainersEmailCheck, t8)

	t9, e9 := STP003.Type(false)
	assert.Nil(t, e9)
	assert.NotNil(t, t9)
	assert.Equal(t, analysisrequest.PypiStaticAnalysisShadyLinks, t9)

	ta, ea := STP010.Type(false)
	assert.Nil(t, ea)
	assert.NotNil(t, ta)
	assert.Equal(t, analysisrequest.PypiStaticAnalysisShadyLinks, ta)

	tb, eb := STN010.Type(false)
	assert.Nil(t, eb)
	assert.NotNil(t, tb)
	assert.Equal(t, analysisrequest.NPMStaticAnalysisShadyLinks, tb)
}
