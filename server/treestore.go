package server

import (
	"github.com/danslimmon/impulse/common"
)

// Treestore provides read and write access to Tree structs persisted to the Datastore.
type Treestore interface {
	GetTree(string) (*common.TreeNode, error)
}

// BasicTreestore is a Treestore implementation in which trees are stored in a basic,
// text-editor-centric serialization format.
//
// The basic format consists of a sequence of lines. A line that is not indented is a direct child
// of the root node. A line that is indented (by any number of consecutive tab characters at the
// beginning of the line) represents a direct child of the next line down with an indentation level
// one less. The bottom line of a tree representation must not be indented.
//
// For examples, see treestore_test.go.
type BasicTreestore struct {
	datastore Datastore
}

// GetTree retrieves the tree with the given name from the persistent Datastore.
func (ts *BasicTreestore) GetTree(treename string) (*common.TreeNode, error) {
	return nil, nil
}

// NewBasicTreestore returns a BasicTreestore with the given underlying datastore.
func NewBasicTreestore(datastore Datastore) *BasicTreestore {
	return &BasicTreestore{datastore: datastore}
}
