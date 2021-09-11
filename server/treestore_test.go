package server

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/danslimmon/impulse/common"
)

func TestBasicTreestore_GetTree(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ds := NewFilesystemDatastore("./testdata")
	ts := NewBasicTreestore(ds)
	rslt, err := ts.GetTree("make_pasta")
	assert.Nil(err)
	//@DEBUG
	makePasta := common.MakePasta()
	t.Logf("MakePasta: %s\n", makePasta.String())
	t.Logf("rslt: %s\n", rslt.String())
	assert.True(rslt.Equal(common.MakePasta()))
}
