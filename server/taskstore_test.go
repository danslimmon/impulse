package server

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/danslimmon/impulse/common"
)

func TestBasicTaskstore_GetList(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ds := NewFilesystemDatastore("./testdata")
	ts := NewBasicTaskstore(ds)
	taskList, err := ts.GetList("make_pasta")
	assert.Nil(err)
	rslt := taskList[0]

	makePasta := common.MakePasta()[0]
	assert.True(rslt.RootNode.Equal(makePasta.RootNode))
}

func TestBasicTaskstore_GetList_Malformed(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ds := NewFilesystemDatastore("./testdata")
	ts := NewBasicTaskstore(ds)

	paths := []string{
		"malformed/zero_length",
		"malformed/excess_delta_indent",
		"malformed/missing",
	}
	for _, p := range paths {
		t.Log(p)
		taskList, err := ts.GetList(p)
		t.Log(err)
		assert.Empty(taskList)
		assert.NotNil(err)
	}
}

func TestBasicTaskstore_PutList(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	a := common.NewTreeNode("alpha")
	a.AddChild(common.NewTreeNode("zulu"))
	a.AddChild(common.NewTreeNode("yankee"))
	b := common.NewTreeNode("bravo")

	tempDir, err := ioutil.TempDir("", "impulse_XXXXXXXXXXXX")
	defer os.RemoveAll(tempDir)
	assert.Nil(err)

	ds := NewFilesystemDatastore(tempDir)
	ts := NewBasicTaskstore(ds)
	err = ts.PutList("foo", []*common.Task{
		common.NewTask(a),
		common.NewTask(b),
	})
	assert.Nil(err)

	taskList, err := ts.GetList("foo")
	assert.Nil(err)
	assert.True(taskList[0].RootNode.Equal(a))
	assert.True(taskList[1].RootNode.Equal(b))
}
