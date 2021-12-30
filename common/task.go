package common

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

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
