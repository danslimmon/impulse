package common

import (
	"sync"
)

type Titler interface {
	GetTitle() string
	SetTitle(string)
}

type TreeNode struct {
	Children []*TreeNode
	Referent Titler
	mu       sync.Mutex
}

func (n *TreeNode) AddChild(childNode *TreeNode) {
	n.mu.Lock()
	n.Children = append(n.Children, childNode)
	n.mu.Unlock()
}

func NewTreeNode(referent Titler) *TreeNode {
	return &TreeNode{
		Children: []*TreeNode{},
		Referent: referent,
	}
}
