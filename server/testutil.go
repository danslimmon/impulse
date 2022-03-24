package server

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"sync"
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

// NewBasicTaskstoreWithTestdata returns a BasicTaskstore based on a clone of the server/testdata
// directory in a tempdir.
//
// NewBasicTaskstoreWithTestdata also returns a function to call when the test is over. Calling this
// function will remove the temporary directory.
func NewBasicTaskstoreWithTestdata() (*BasicTaskstore, func()) {
	ds, cleanup := newFSDatastoreWithTestdata()
	return NewBasicTaskstore(ds), cleanup
}

// ap is a singleton addrPool that we use to provision addrs for tests to listen on.
var ap *addrPool

// addrPool provisions unique addrs for tests to listen on.
//
// An addr is a string of the form "<IP>:<port>". We will provision up to 10000 unique addrs. If
// more than 10000 are requested, addrPool.Get panics.
type addrPool struct {
	port int
	mu   sync.Mutex
}

func (a *addrPool) Get() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.port == 0 {
		a.port = 30000
	} else if a.port >= 39999 {
		panic("more than 10000 addrs requested from addrPool")
	} else {
		a.port = a.port + 1
	}
	return fmt.Sprintf("127.0.0.1:%d", a.port)
}

// listenAddr returns a loopback address for a test instance of the API to listen on.
//
// The address is guaranteed not to be in use by any other test in the suite.
func listenAddr() string {
	if ap == nil {
		ap = new(addrPool)
	}
	return ap.Get()
}

// NewServerWithTestdata returns a Server based on a clone of the server/testdata directory in a
// teempdir.
//
// NewServerWithTestdata also returns a function to call when the test is over. Calling this
// function will remove the temporary directory.
func NewServerWithTestdata() (*Server, func()) {
	ts, cleanup := NewBasicTaskstoreWithTestdata()
	s := &Server{
		taskstore: ts,
	}
	addr := listenAddr()
	err := s.Start(addr)
	if err != nil {
		panic(err.Error())
	}

	return s, cleanup
}
