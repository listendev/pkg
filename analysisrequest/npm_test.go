package analysisrequest

import (
	"errors"
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
	aaa, err := NewNPM(NPMInstallWhileDynamicInstrumentation, id, prio, force, name, vers, shasum)
	assert.Nil(t, err)
	assert.NotNil(t, aaa)

	arn, ok := aaa.(*NPM)
	assert.True(t, ok)
	assert.NotNil(t, arn)

	enrichWithAI, err := arn.Switch(NPMInstallWhileDynamicInstrumentationAIEnriched)
	assert.Nil(t, err)
	assert.NotNil(t, enrichWithAI)
	assert.Equal(t, NPMInstallWhileDynamicInstrumentationAIEnriched, enrichWithAI.Type())
	assert.Equal(t, force, enrichWithAI.MustProcess())
	assert.Equal(t, prio, enrichWithAI.Prio())

	_, noEcoErr := enrichWithAI.(*NPM).Switch(Nop)
	if assert.Error(t, noEcoErr) {
		assert.Equal(t, "couldn't switch the current NPM analysis request to an analysis request with a type without ecosystem", noEcoErr.Error())
	}
}

func TestSetPrio(t *testing.T) {
	id := "1524854487523524608"
	prio := uint8(5)
	force := false
	name := "chalk"
	vers := "5.2.0"
	shasum := "249623b7d66869c673699fb66d65723e54dfcfb3"
	aaa, err := NewNPM(NPMInstallWhileDynamicInstrumentation, id, prio, force, name, vers, shasum)
	assert.Nil(t, err)
	assert.NotNil(t, aaa)

	arn, ok := aaa.(*NPM)
	assert.True(t, ok)
	assert.NotNil(t, arn)
	assert.Equal(t, prio, arn.Prio())

	arn.SetPrio(uint8(2))
	assert.Equal(t, uint8(2), arn.Prio())
}

func TestErrors(t *testing.T) {
	assert.True(t, errors.As(ErrGivenVersionNotFoundOnNPM, &NPMFillError{}))
	assert.True(t, errors.As(ErrGivenShasumDoesntMatchGivenVersionOnNPM, &NPMFillError{}))
}
