package analysisrequest

import (
	"testing"

	"github.com/listendev/pkg/ecosystem"
	"github.com/stretchr/testify/assert"
)

func TestGetEcosystemFrom(t *testing.T) {
	_, err0 := ecosystem.FromString("nope")
	if assert.NotNil(t, err0) {
		assert.Error(t, err0)
	}

	// NPM
	eco1, err1 := ecosystem.FromString("NPM")
	assert.Nil(t, err1)
	assert.NotNil(t, eco1)

	eco2, err2 := ecosystem.FromString("npm")
	assert.Nil(t, err2)
	assert.NotNil(t, eco2)

	assert.Equal(t, eco1, eco2)

	// PyPi
	eco3, err3 := ecosystem.FromString("PyPi")
	assert.Nil(t, err3)
	assert.NotNil(t, eco3)

	eco4, err4 := ecosystem.FromString("pypi")
	assert.Nil(t, err4)
	assert.NotNil(t, eco4)

	assert.Equal(t, eco3, eco4)
}
