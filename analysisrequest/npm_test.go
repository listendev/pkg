package analysisrequest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSwitch(t *testing.T) {
	id := "1524854487523524608"
	prio := uint8(5)
	force := false
	name := "chalk"
	vers := "5.2.0"
	shasum := "249623b7d66869c673699fb66d65723e54dfcfb3"
	aaa, err := NewNPM(NPMInstallWhileFalco, id, prio, force, name, vers, shasum)
	assert.Nil(t, err)
	assert.NotNil(t, aaa)

	arn, ok := aaa.(*NPM)
	assert.True(t, ok)
	assert.NotNil(t, arn)

	enrichWithGPT, err := arn.Switch(NPMGPT4InstallWhileFalco)
	assert.Nil(t, err)
	assert.NotNil(t, enrichWithGPT)
	assert.Equal(t, NPMGPT4InstallWhileFalco, enrichWithGPT.Type())
	assert.Equal(t, force, enrichWithGPT.MustProcess())
	assert.Equal(t, prio, enrichWithGPT.Prio())

	_, noEcoErr := enrichWithGPT.(*NPM).Switch(Nop)
	if assert.Error(t, noEcoErr) {
		assert.Equal(t, "couldn't switch the current NPM analysis request to an analysis request with a type without ecosystem", noEcoErr.Error())
	}
}