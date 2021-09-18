package common

import (
	"errors"
	"fmt"
	"strings"
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
	Parent   *TreeNode
	Children []*TreeNode
	Referent string
	mu       sync.Mutex
}

// Depth returns the number of ancestors of n.
func (n *TreeNode) Depth() int {
	i := 0
	for n.Parent != nil {
		i++
		n = n.Parent
	}
	return i
}

// AddChild makes childNode a child of n.
//
// childNode is added at the end of n.Children.
func (n *TreeNode) AddChild(childNode *TreeNode) {
	n.mu.Lock()
	defer n.mu.Unlock()
	childNode.Parent = n
	n.Children = append(n.Children, childNode)
}

// InsertChild makes childNode a child of n, placing it at the given position among n.Children.
//
// ind is the index in n.Children at which childNode will be inserted.
func (n *TreeNode) InsertChild(ind int, childNode *TreeNode) {
	n.mu.Lock()
	defer n.mu.Unlock()

	childNode.Parent = n
	n.Children = append(n.Children, nil)
	copy(n.Children[ind+1:], n.Children[ind:])
	n.Children[ind] = childNode
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

// WalkFromTop walks the tree rooted at n, calling fn for each TreeNode, including n.
//
// All errors that arise are filtered by fn: see the TreeWalkFunc documentation for details.
//
// WalkFromTop visits nodes from the top to the bottom of the task. For example, for the "make
// pasta" task (server/testdata/make_pasta), fn would be called on the nodes in this order:
//
// - put water in pot
// - put pot on burner
// - turn burner on
// - boil water
// - put pasta in water
// - [b cooked]
// - drain pasta
// - make pasta
//
// Since this is a depth-first walk, SkipSubtree is not treated specially: if fn returns
// SkipSubtree, the walk exits with that error.
func (n *TreeNode) WalkFromTop(fn TreeWalkFunc) error {
	var err error
	for _, cn := range n.Children {
		err = cn.WalkFromTop(fn)
		if err != nil {
			return err
		}
	}

	err = fn(n)
	if err != nil {
		return err
	}

	return nil
}

// String returns a string representation of n, for use in logging and debugging.
//
// One should not use the return value of String() to compare TreeNodes. Instead, one should use
// Equal() and/or write a custom TreeWalkFunc. String() is only for convenience.
func (n *TreeNode) String() string {
	var b strings.Builder
	n.Walk(func(m *TreeNode) error {
		b.WriteString(fmt.Sprintf("%s%s\n", strings.Repeat("\t", m.Depth()), m.Referent))
		return nil
	})
	return b.String()
}

// Equal determines whether a is identical to b.
//
// Equal walks both trees and compares the corresponding nodes in a and b. If the two nodes' depth
// and referent are equal, then the nodes are equal.
func (a *TreeNode) Equal(b *TreeNode) bool {
	if a == nil || b == nil {
		// technically nil should equal nil, but… i don't want a situation where anybody's ever
		// passing nil to this function
		panic("called *TreeNode.Equal on nil TreeNode")
	}

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

// NewTreeNode returns a TreeNode with the given Referent.
func NewTreeNode(referent string) *TreeNode {
	return &TreeNode{
		Children: []*TreeNode{},
		Referent: referent,
	}
}
