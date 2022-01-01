package server

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/danslimmon/impulse/common"
)

// newBasicTaskstoreWithTestdata returns a BasicTaskstore based on a clone of the server/testdata
// directory in a tempdir.
//
// newBasicTaskstoreWithTestdata also returns a function to call when the test is over. Calling this
// function will remove the temporary directory.
func newBasicTaskstoreWithTestdata() (*BasicTaskstore, func()) {
	ds, cleanup := newFSDatastoreWithTestdata()
	return NewBasicTaskstore(ds), cleanup
}

func TestBasicTaskstore_derefLineId(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	type testCase struct {
		LineId      common.LineID
		ExpListName string
		ExpLineNo   int
		ExpErr      bool
	}

	testCases := []testCase{
		// happy path
		testCase{
			LineId:      common.GetLineID("make_pasta", "\tboil water"),
			ExpListName: "make_pasta",
			ExpLineNo:   3,
			ExpErr:      false,
		},
		testCase{
			LineId:      common.GetLineID("multiple_nested", "\t\tsubsubtask 0"),
			ExpListName: "multiple_nested",
			ExpLineNo:   0,
			ExpErr:      false,
		},
		testCase{
			LineId:      common.GetLineID("multiple_nested", "task 1"),
			ExpListName: "multiple_nested",
			ExpLineNo:   4,
			ExpErr:      false,
		},
		testCase{
			LineId:      common.LineID("multiple_nested:0"),
			ExpListName: "multiple_nested",
			ExpLineNo:   0,
			ExpErr:      false,
		},
		// sad path
		testCase{
			LineId: common.LineID("malformed line ID (no slash)"),
			ExpErr: true,
		},
		testCase{
			LineId: common.GetLineID("no_such_file", "blah blah"),
			ExpErr: true,
		},
		testCase{
			LineId: common.GetLineID("make_pasta", "line that doesn't exist"),
			ExpErr: true,
		},
	}

	for _, tc := range testCases {
		ds := NewFilesystemDatastore("./testdata")
		ts := NewBasicTaskstore(ds)
		rsltListName, rsltLineNo, rsltErr := ts.derefLineId(tc.LineId)
		if !tc.ExpErr {
			assert.Nil(rsltErr)
			assert.Equal(tc.ExpListName, rsltListName)
			assert.Equal(tc.ExpLineNo, rsltLineNo)
		} else {
			assert.NotNil(rsltErr)
		}
	}
}

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

	ts, cleanup := newBasicTaskstoreWithTestdata()
	defer cleanup()

	err := ts.PutList("foo", []*common.Task{
		common.NewTask(a),
		common.NewTask(b),
	})
	assert.Nil(err)

	taskList, err := ts.GetList("foo")
	assert.Nil(err)
	assert.True(taskList[0].RootNode.Equal(a))
	assert.True(taskList[1].RootNode.Equal(b))
}

func TestBasicTaskstore_InsertTask(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	a := common.NewTreeNode("alpha")
	a.AddChild(common.NewTreeNode("zulu"))
	a.AddChild(common.NewTreeNode("yankee"))
	b := common.NewTreeNode("bravo")

	ts, cleanup := newBasicTaskstoreWithTestdata()
	defer cleanup()

	// Insert alpha at the top of make_pasta
	err := ts.InsertTask(common.LineID("make_pasta:0"), common.NewTask(a))
	assert.Nil(err)

	taskList, err := ts.GetList("make_pasta")
	assert.Nil(err)
	assert.True(a.Equal(taskList[0].RootNode))

	// Insert bravo at the bottom of multiple_nested
	err = ts.InsertTask(common.GetLineID("multiple_nested", "task 1"), common.NewTask(b))
	assert.Nil(err)

	taskList, err = ts.GetList("multiple_nested")
	assert.Nil(err)
	assert.True(b.Equal(taskList[2].RootNode))
}

// Tests that ArchiveLine works when given an ID that corresponds to a subtask.
func TestBasicTaskstore_ArchiveLine_Subtask(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ds, cleanup := newFSDatastoreWithTestdata()
	ts := NewBasicTaskstore(ds)
	defer cleanup()

	err := ts.ArchiveLine(common.GetLineID("make_pasta", "\t\tput water in pot"))
	assert.Nil(err)

	// make sure that the archive operation didn't cause malformation of the list file
	_, err = ts.GetList("make_pasta")
	assert.Nil(err)

	// make sure that the line actually got removed from the list file
	b, err := ds.Get("make_pasta")
	assert.Nil(err)
	assert.False(regexp.MustCompile("put water in pot").Match(b))

	// make sure that the history file now contains the line we archived
	b, err = ds.Get("history")
	assert.Nil(err)
	assert.True(regexp.MustCompile("^[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9] \t\tput water in pot$").Match(b))
}

// Tests that ArchiveLine works when given an ID that corresponds to a task.
func TestBasicTaskstore_ArchiveLine_Task(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ds, cleanup := newFSDatastoreWithTestdata()
	ts := NewBasicTaskstore(ds)
	defer cleanup()

	err := ts.ArchiveLine(common.GetLineID("make_pasta", "make pasta"))
	assert.Nil(err)

	// make sure that the archive operation didn't cause malformation of the list file
	_, err = ts.GetList("make_pasta")
	assert.Nil(err)

	// make sure that the line actually got removed from the list file
	b, err := ds.Get("make_pasta")
	assert.Nil(err)
	assert.False(regexp.MustCompile("make pasta").Match(b))

	// make sure that the history file now contains the line we archived
	b, err = ds.Get("history")
	assert.Nil(err)
	assert.True(regexp.MustCompile("^[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9] make pasta$").Match(b))
}
