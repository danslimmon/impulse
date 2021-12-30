package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLineID(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// trailing newlines ignored
	assert.Equal(
		GetLineID("foo", "\t\tbar"),
		GetLineID("foo", "\t\tbar\n"),
	)
	assert.Equal(
		GetLineID("foo", "\t\tbar"),
		GetLineID("foo", "\t\tbar\n\n\n"),
	)

	// no collisions
	assert.NotEqual(
		GetLineID("foo", "bar"),
		GetLineID("not_foo", "bar"),
	)
	assert.NotEqual(
		GetLineID("foo", "bar"),
		GetLineID("foo", "not_bar"),
	)
}
