package server

import (
	"bytes"
	"fmt"

	"github.com/danslimmon/impulse/common"
)

/*
// Taskstore provides read and write access to Tree structs persisted to the Datastore.
type Taskstore interface {
	Get(string) ([]*common.Task, error)
	Put(string, []*common.Task) error
}
*/

// BasicTaskstore is a Taskstore implementation in which trees are stored in a basic,
// text-editor-centric serialization format.
//
// The basic format consists of a sequence of lines. A line that is not indented is a direct child
// of the root node. A line that is indented (by any number of consecutive tab characters at the
// beginning of the line) represents a direct child of the next line down with an indentation level
// one less. The bottom line of a tree representation must not be indented.
//
// For examples, see treestore_test.go.
type BasicTaskstore struct {
	datastore Datastore
}

// parseLine parses a line of basic-format tree data.
//
// It returns the integer number of tabs that occur at the beginning of the line (its indent level)
// and the remaining text of the line as a string.
func (ts *BasicTaskstore) parseLine(line []byte) (int, string) {
	textBytes := bytes.TrimLeft(line, "\t")
	indent := len(line) - len(textBytes)
	return indent, string(textBytes)
}

// Get retrieves the task list with the given name from the persistent Datastore.
func (ts *BasicTaskstore) Get(name string) ([]*common.Task, error) {
	b, err := ts.datastore.Get(name)
	if err != nil {
		return nil, err
	}

	// lines ends up just being splut in reverse. so lines is all the lines in the file, from the
	// bottom to the top of the file. that's how we want it for constructing the tree further down.
	splut := bytes.Split(b, []byte("\n"))
	if len(splut) == 0 {
		return []*common.Task{}, nil
	}

	nLines := len(splut)
	lines := make([][]byte, nLines)
	for i := range splut {
		lines[nLines-1-i] = splut[i]
	}

	if len(lines[0]) == 0 {
		lines = lines[1:]
		nLines = nLines - 1
	} else {
		return nil, fmt.Errorf("error: data for tree '%s' does not end in newline", name)
	}

	rootNode := common.NewTreeNode("")
	// prevIndent is the indent of the previous line. remember: lines contains all the file's lines
	// _from the bottom to the top!_
	//
	// the indent of the "root node", we can say, is -1. such that the bottommost line, which should
	// start with 0 indent, is a child of that root node.
	prevIndent := -1
	// prevNode points to the line parsed in the previous iteration of the loop.
	prevNode := rootNode
	for i, line := range lines {
		indent, text := ts.parseLine(line)
		deltaIndent := indent - prevIndent
		newNode := common.NewTreeNode(text)

		if deltaIndent == 1 {
			// this is a child of the previous node
			prevNode.AddChild(newNode)
		} else if deltaIndent == 0 {
			// this is a sibling of the previous node. it goes _before_ the previous node parsed.
			prevNode.Parent.InsertChild(0, newNode)
		} else if deltaIndent < 0 {
			// we've gone back up to an ancestor node. figure out which one and add the child there
			// (again, before the previous node parsed)
			//
			// ascend the tree by however much it takes to get back to the ancestor of node that
			// newNode is a child of
			ancestorNode := prevNode
			for i := 0; i < -deltaIndent-1; i++ {
				ancestorNode = ancestorNode.Parent
			}
			ancestorNode.InsertChild(0, newNode)
		} else {
			return nil, fmt.Errorf(
				"error parsing line %d of tree '%s': unexpected deltaIndent = %d",
				i,
				name,
				deltaIndent,
			)
		}
		prevNode = newNode
		prevIndent = indent
	}

	// Now we take all the children of rootNode and load them into the list to be returned.
	rslt := make([]*common.Task, 0)
	for _, n := range rootNode.Children {
		rslt = append(rslt, common.NewTask(n))
	}
	return rslt, nil
}

// Put writes taskList to the Datastore as name.
func (ts *BasicTaskstore) Put(name string, taskList []*common.Task) error {
	b := []byte{}
	for _, t := range taskList {
		t.RootNode.WalkFromTop(func(n *common.TreeNode) error {
			b = append(b, bytes.Repeat([]byte("\t"), n.Depth())...)
			b = append(b, []byte(n.Referent)...)
			b = append(b, []byte("\n")...)
			return nil
		})
	}
	return ts.datastore.Put(name, b)
}

// NewBasicTaskstore returns a BasicTaskstore with the given underlying datastore.
func NewBasicTaskstore(datastore Datastore) *BasicTaskstore {
	return &BasicTaskstore{datastore: datastore}
}
