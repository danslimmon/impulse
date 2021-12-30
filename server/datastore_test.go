package server

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

// cloneTestData copies server/testdata to a new tempdir and returns that tempdir's path.
//
// cloneTestData also returns a function to call when the test is over. Calling this function will
// remove the temporary directory.
func cloneTestData() (string, func()) {
	tempDir, err := ioutil.TempDir("", "impulse_*")
	if err != nil {
		panic("unable to create tempdir for testdata clone: " + err.Error())
	}
	if err := exec.Command("cp", "-rp", "./testdata/", tempDir+"/").Run(); err != nil {
		panic("unable to clone testdata to tempdir: " + err.Error())
	}
	return tempDir, func() {
		os.RemoveAll(tempDir)
	}
}

// newFSDatastoreWithTestdata returns a FilesystemDatastore whose filesystem has been cloned from
// server/testdata.
//
// newFSDatastoreWithTestdata also returns a function to call when the test is over. Calling this
// function will remove the temporary directory that is the FilesystemDatastore's rootDir.
func newFSDatastoreWithTestdata() (*FilesystemDatastore, func()) {
	tempDir, cleanup := cloneTestData()
	return NewFilesystemDatastore(tempDir), cleanup
}

func TestFilesystemDatastore_Append(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ds, cleanup := newFSDatastoreWithTestdata()
	defer cleanup()

	// Append to a file that doesn't exist yet
	err := ds.Append("foo", []byte("first line\n"))
	assert.Nil(err)
	rslt, err := ds.Get("foo")
	assert.Nil(err)
	assert.Equal([]byte("first line\n"), rslt, fmt.Sprintf("unexpected file contents: '%s'", string(rslt)))

	// Append to a file that already exists
	err = ds.Append("foo", []byte("second line\n"))
	assert.Nil(err)
	rslt, err = ds.Get("foo")
	assert.Nil(err)
	assert.Equal([]byte("first line\nsecond line\n"), rslt, fmt.Sprintf("unexpected file contents: '%s'", string(rslt)))
}
