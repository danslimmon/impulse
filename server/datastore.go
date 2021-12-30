package server

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const DataDir = "/Users/danslimmon/j_workspace/impulse"

// Datastore is an interface to raw marshaled Impulse data.
//
// A Datastore implementation is responsible for:
//
//   - Mapping task list names to the corresponding objects in the database (e.g. files on disk)
//   - Reading and writing the data in those objects
type Datastore interface {
	// Get returns the contents of the file corresponding to the given task list
	Get(string) ([]byte, error)
	// Put writes the given data to the file with the given name.
	Put(string, []byte) error
	// Append appends to the file identified by the given name, the given bytes.
	//
	// If no file with the given name exists, Append creates the file and sets its contents to the
	// provided []byte.
	//
	// The caller is responsible for including a \n.
	Append(string, []byte) error
}

// FilesystemDatastore is a Datastore implementation in which trees are marshaled into files in a
// straightforward directory tree according to their names.
//
// The root of the filesystem in which data is stored is RootDir.
//
// For example, if RootDir is "/path", the contents of the tree with name "pers" can be found in the
// file "/path/pers". If a given tree name contains slashes, they are treated as path separators by
// FilesystemDatastore.
type FilesystemDatastore struct {
	rootDir string
}

// absPath returns the full path to the file that should contain the given task list's marshaled data.
func (ds *FilesystemDatastore) absPath(name string) string {
	return filepath.Join(ds.rootDir, name)
}

// See Datastore interface
func (ds *FilesystemDatastore) Get(name string) ([]byte, error) {
	return ioutil.ReadFile(ds.absPath(name))
}

// See Datastore interface
func (ds *FilesystemDatastore) Put(name string, b []byte) error {
	return ioutil.WriteFile(ds.absPath(name), b, 0644)
}

// See Datastore interface
func (ds *FilesystemDatastore) Append(name string, b []byte) error {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(ds.absPath(name), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err := f.Write(b); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func NewFilesystemDatastore(rootDir string) *FilesystemDatastore {
	return &FilesystemDatastore{rootDir: rootDir}
}
