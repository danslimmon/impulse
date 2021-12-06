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

	type testCase struct {
		ListName string
		Exp      []*common.Task
	}

	testCases := []testCase{
		testCase{
			ListName: "make_pasta",
			Exp:      common.MakePasta(),
		},
		testCase{
			ListName: "multiple_nested",
			Exp:      common.MultipleNested(),
		},
	}

	for _, tc := range testCases {
		ds := NewFilesystemDatastore("./testdata")
		ts := NewBasicTaskstore(ds)
		rslt, err := ts.GetList(tc.ListName)
		assert.Nil(err)
		assert.Equal(len(tc.Exp), len(rslt))
		for i := range tc.Exp {
			if i > len(rslt)-1 {
				break
			}
			assert.Equal(tc.Exp[i], rslt[i])
		}
	}
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
