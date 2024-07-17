package config

import (
	"strings"
	"testing"

	"gotest.tools/assert"
)

func Test_EqualFold(t *testing.T) {
	l, r := "Go=123", "go=456"
	i := strings.Index(l, "=")
	assert.Equal(t, true, strings.EqualFold(l[:i+1], r[:i+1]))
}
