package common

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// LineID represents a line in the data file.
//
// This is an abstraction leak. API clients shouldn't have to know about "lines", since those are an
// implementation detail of the storage backend. But, it gets us off the ground for now.
//
// A Line ID is composed of two parts, separated by a colon. The first part is the name of a task
// list, e.g. `make_pasta`. The second part is either 0 (which indicates the top line of the file),
// or a SHA256 sum of the line's content, indentation included and trailing whitespace excluded (see
// GetLineID).
type LineID string

// GetLineID returns the line ID for the line in the list identified by listName, with content s.
//
// s will be stripped of trailing newlines before the ID is calculated.
func GetLineID(listName, s string) LineID {
	s = strings.TrimRight(s, "\n")
	shasumArray := sha256.Sum256([]byte(s))
	shasum := shasumArray[:]
	digest := hex.EncodeToString(shasum)
	return LineID(fmt.Sprintf("%s:%s", listName, digest))
}

type Task struct {
	RootNode *TreeNode `json:"tree"`
}

// Equal determines whether the tasks a and b are equal.
func (a *Task) Equal(b *Task) bool {
	return a.RootNode.Equal(b.RootNode)
}

func NewTask(n *TreeNode) *Task {
	return &Task{RootNode: n}
}
