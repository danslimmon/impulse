package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TreeNode is a node in a Tree.
type TreeNode struct {
	task     *Task
	parent   *TreeNode
	children []*TreeNode
}

func rootTreeNode() *TreeNode {
	return &TreeNode{
		task:     nil,
		parent:   nil,
		children: []*TreeNode{},
	}
}

// Tree is the tree of tasks in Impulse.
//
// It's largely adapted from tview.TreeView.
type Tree struct {
	// pull in default functions to satisfy the tview.Primitive interface
	*tview.Box

	root *TreeNode
}

// NewTree returns a new Tree with a root node.
func NewTree() *Tree {
	return &Tree{
		Box:  tview.NewBox(),
		root: rootTreeNode(),
	}
}

// Draw draws this primitive onto the screen.
func (t *Tree) Draw(screen tcell.Screen) {
	t.Box.DrawForSubclass(screen, t)
	if t.root == nil {
		return
	}
	screenWidth, _ := screen.Size()

	tview.Print(screen, "boil water", 4, 0, screenWidth, tview.AlignLeft, tcell.ColorReset)
	tview.Print(screen, "make pasta", 0, 1, screenWidth, tview.AlignLeft, tcell.ColorBlue)
}

// this is the one we have to implement
func (tree *Tree) InputHandler() func(*tcell.EventKey, func(p tview.Primitive)) {
	return func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {}
}
