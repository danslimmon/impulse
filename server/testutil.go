package server

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
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

	cwd, err := os.Getwd()
	if err != nil {
		panic("unable to get current working directory: " + err.Error())
	}

	var srcDir string
	if path.Base(cwd) == "impulse" {
		srcDir = "testdata"
	} else {
		srcDir = "../testdata"
	}

	if err := exec.Command("cp", "-rp", srcDir+"/", tempDir+"/").Run(); err != nil {
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

// newBasicTaskstoreWithTestdata returns a BasicTaskstore based on a clone of the server/testdata
// directory in a tempdir.
//
// newBasicTaskstoreWithTestdata also returns a function to call when the test is over. Calling this
// function will remove the temporary directory.
func NewBasicTaskstoreWithTestdata() (*BasicTaskstore, func()) {
	ds, cleanup := newFSDatastoreWithTestdata()
	return NewBasicTaskstore(ds), cleanup
}
