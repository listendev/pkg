package npm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitOrgName(t *testing.T) {
	org, pkg := SplitName("@vue/devtools")
	assert.Equal(t, "@vue", org)
	assert.Equal(t, "devtools", pkg)
}

func TestSplitName(t *testing.T) {
	org, pkg := SplitName("vu3")
	assert.Equal(t, "", org)
	assert.Equal(t, "vu3", pkg)
}
