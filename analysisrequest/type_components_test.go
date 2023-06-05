package analysisrequest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEcosystemFrom(t *testing.T) {
	eco1, err1 := GetEcosystemFrom("NPM")
	assert.Nil(t, err1)
	assert.NotNil(t, eco1)

	eco2, err2 := GetEcosystemFrom("npm")
	assert.Nil(t, err2)
	assert.NotNil(t, eco2)

	assert.Equal(t, eco1, eco2)

	_, err3 := GetEcosystemFrom("nope")
	assert.NotNil(t, err3)
	assert.Error(t, err3)
}
