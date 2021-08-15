package server

import (
	"io/ioutil"
	"path/filepath"
)

// Datastore is an interface to raw marshaled Impulse data.
//
// A Datastore implementation is responsible for:
//
//   - Mapping tree names to the corresponding objects in the database (e.g. files on disk)
//   - Reading and writing the data in those objects
type Datastore interface {
	// GetTreeData returns the contents of the file corresponding to the given tree.
	GetTreeData(string) ([]byte, error)
}

// FilesystemDatastore is a Datastore implementation in which trees are marshaled onto files in a
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

// absPath returns the full path to the file that should contain the given tree's marshaled data.
func (ds *FilesystemDatastore) absPath(treename string) string {
	return filepath.Join(ds.rootDir, treename)
}

func (ds *FilesystemDatastore) GetTreeData(treename string) ([]byte, error) {
	return ioutil.ReadFile(ds.absPath(treename))
}
