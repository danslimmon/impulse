package server

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/danslimmon/impulse/common"
)

// Taskstore provides read and write access to Tree structs persisted to the Datastore.
type Taskstore interface {
	GetList(string) ([]*common.Task, error)
	PutList(string, []*common.Task) error
	ArchiveLine(common.LineID) error
}

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

// derefLineId takes a line ID and determines the corresponding list name and line number.
//
// The line number returned is zero-indexed (the first line of the file is line 0).
func (ts *BasicTaskstore) derefLineId(lineId common.LineID) (string, int, error) {
	parts := strings.SplitN(string(lineId), ":", 2)
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("malformatted line ID `%s`", string(lineId))
	}

	listName := parts[0]
	b, err := ts.datastore.Get(listName)
	if err != nil {
		return "", 0, err
	}

	lines := bytes.Split(b, []byte("\n"))
	lineNo := -1
	for n, line := range lines {
		if common.GetLineID(listName, string(line)) == lineId {
			lineNo = n
		}
	}
	if lineNo == -1 {
		return "", 0, fmt.Errorf("no line with ID `%s`", string(lineId))
	}

	return listName, lineNo, nil
}

// Get retrieves the task list with the given name from the persistent Datastore.
func (ts *BasicTaskstore) GetList(name string) ([]*common.Task, error) {
	b, err := ts.datastore.Get(name)
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		return nil, fmt.Errorf("data for tree '%s' is zero-length", name)
	}

	// lines ends up just being splut in reverse. so lines is all the lines in the file, from the
	// bottom to the top of the file. that's how we want it for constructing the tree further down.
	splut := bytes.Split(b, []byte("\n"))
	if len(splut) < 2 {
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
		return nil, fmt.Errorf("data for tree '%s' does not end in newline", name)
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
		} else if deltaIndent <= 0 {
			// we've gone back up to an ancestor node. figure out which one and add the child there
			// (again, before the previous node parsed)
			//
			// ascend the tree by however much it takes to get back to the ancestor of node that
			// newNode is a child of
			ancestorNode := prevNode.Parent
			for ancestorNode.Depth() != prevNode.Depth()+deltaIndent-1 {
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
		n.Parent = nil
		rslt = append(rslt, common.NewTask(n))
	}
	return rslt, nil
}

// Put writes taskList to the Datastore as name.
func (ts *BasicTaskstore) PutList(name string, taskList []*common.Task) error {
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

// historyLine returns a line for the history file based on the given line from an impulse file.
//
// History lines are of the form:
//
//     2021-12-30T19:24:48 [full contents of b, including any leading whitespace]
func (ts *BasicTaskstore) historyLine(b []byte) []byte {
	now := time.Now()
	// make sure this is UTC before using it ^
	timestamp := now.Format("2006-01-02T15:04:05")
	return []byte(fmt.Sprintf("%s %s", timestamp, b))
}

// ArchiveLine archives the line identified by lineId.
//
// lineId may refer either to a subtask or a task proper.
func (ts *BasicTaskstore) ArchiveLine(lineId common.LineID) error {
	listName, lineNo, err := ts.derefLineId(lineId)
	if err != nil {
		return err
	}

	b, err := ts.datastore.Get(listName)
	if err != nil {
		return err
	}

	lines := bytes.Split(b, []byte("\n"))
	// Will panic on index-out-of-range, but it should. That means there's a bug in derefLineId.
	removedLine := make([]byte, len(lines[lineNo]))
	copy(removedLine, lines[lineNo])

	lines = append(lines[0:lineNo], lines[lineNo+1:]...)
	// Add empty line so that file will end with a newline
	lines = append(lines, []byte{})
	b = bytes.Join(lines, []byte("\n"))

	ts.datastore.Append("history", ts.historyLine(removedLine))

	return ts.datastore.Put(listName, b)
}

// NewBasicTaskstore returns a BasicTaskstore with the given underlying datastore.
func NewBasicTaskstore(datastore Datastore) *BasicTaskstore {
	return &BasicTaskstore{datastore: datastore}
}
