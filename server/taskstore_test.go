package server

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/danslimmon/impulse/common"
)

func TestBasicTaskstore_Get(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ds := NewFilesystemDatastore("./testdata")
	ts := NewBasicTaskstore(ds)
	taskList, err := ts.Get("make_pasta")
	rslt := taskList[0]
	assert.Nil(err)
	//@DEBUG
	makePasta := common.MakePasta()[0]
	t.Logf("MakePasta: %s\n", makePasta.RootNode.String())
	t.Logf("rslt: %s\n", rslt.RootNode.String())
	assert.True(rslt.RootNode.Equal(makePasta.RootNode))
}
