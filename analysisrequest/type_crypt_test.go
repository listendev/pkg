package analysisrequest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	input := "22"
	a, e := NPMDepsDev.Encrypt(input)
	assert.Nil(t, e)
	assert.NotNil(t, a)

	res, err := Decrypt(a)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, fmt.Sprintf("%s@%s", NPMDepsDev.String(), input), string(res))
}
