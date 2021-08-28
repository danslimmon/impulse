package common

import (
	"errors"
	"sync"
)

// SkipSubtree is returned by a TreeWalkFunc in order to skip descending into a given subtree
var SkipSubtree = errors.New("skip this subtree")

// TreeWalkFunc is a function called on each node in a tree by TreeNode.Walk.
//
// If TreeWalkFunc returns an error, and that error is not SkipSubtree, the rest of the walk is
// aborted. If SkipSubtree is returned, Walk will not descend into node's children, but will proceed
// with the rest of the walk as normal.
type TreeWalkFunc func(*TreeNode) error

type TreeNode struct {
	Children []*TreeNode
	Referent string
	mu       sync.Mutex
}

func (n *TreeNode) AddChild(childNode *TreeNode) {
	n.mu.Lock()
	n.Children = append(n.Children, childNode)
	n.mu.Unlock()
}

// Walk walks the tree rooted at n, calling fn for each TreeNode, including n.
//
// All errors that arise are filtered by fn: see the TreeWalkFunc documentation for details.
//
// After a node is visited, each of its children is walked, in the order in which they would be
// executed. For example, for the "make pasta" task (server/testdata/make_pasta), fn would be called
// on the nodes in this order:
//
// - make pasta
// - boil water
// - put water in pot
// - put pot on burner
// - turn burner on
// - put pasta in water
// - [b cooked]
// - drain pasta
func (n *TreeNode) Walk(fn TreeWalkFunc) error {
	err := fn(n)
	if err == SkipSubtree {
		return nil
	}
	if err != nil {
		return err
	}

	for _, cn := range n.Children {
		err = cn.Walk(fn)
		if err != nil {
			return err
		}
	}

	return nil
}

// Equal determines whether a is identical to b.
//
// Equal walks both trees and compares the corresponding nodes in a and b. If the two nodes' depth
// and referent are equal, then the nodes are equal.
func (a *TreeNode) Equal(b *TreeNode) bool {
	toChan := func(rootNode *TreeNode, ch chan *TreeNode) {
		rootNode.Walk(func(n *TreeNode) error {
			ch <- n
			return nil
		})
		close(ch)
	}

	aCh := make(chan *TreeNode)
	bCh := make(chan *TreeNode)
	go toChan(a, aCh)
	go toChan(b, bCh)

	for aNode := range aCh {
		bNode := <-bCh
		if len(aNode.Children) != len(bNode.Children) {
			return false
		}
		if aNode.Referent != bNode.Referent {
			return false
		}
	}

	// Make sure b is finished
	_, ok := <-bCh
	if ok {
		return false
	}

	return true
}

func NewTreeNode(referent string) *TreeNode {
	return &TreeNode{
		Children: []*TreeNode{},
		Referent: referent,
	}
}
