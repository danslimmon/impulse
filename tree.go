package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	nodeIndent = "    "
)

// TreeNode is a node in a Tree.
type TreeNode struct {
	// Ordered from bottom to top
	children []*TreeNode
	// nil for the root node
	parent *TreeNode
	// Whether the TreeNode is currently selected in the UI
	selected bool
	task     *Task
}

// Depth returns the number of ancestors that node has.
//
// For example, the Depth of the root node is 0. The depth of one of the root node's children is 1.
func (node *TreeNode) Depth() int {
	if node.parent == nil {
		return 0
	}
	return 1 + node.parent.Depth()
}

// PushTask adds a new TreeNode for the given task.
//
// The new TreeNode, which is returned, will be the topmost child of node.
func (node *TreeNode) PushTask(task *Task) *TreeNode {
	newNode := newTreeNode(node, task)
	node.children = append(node.children, newNode)
	return newNode
}

// Walk traverses this node's subtree in depth-first, bottom-to-top-of-stack order, and calls the
// provided callback function on each traversed node (which includes this node) with the traversed
// node and its parent node (nil for this node).  The callback returns whether traversal should
// continue with the traversed node's child nodes (true) or not recurse any deeper (false).
//
// The return value of Walk is equal to the return value of callback when called on node.
func (node *TreeNode) Walk(callback func(*TreeNode, *TreeNode) bool) bool {
	cont := callback(node, node.parent)
	if !cont {
		return false
	}
	for _, child := range node.children {
		child.Walk(callback)
	}
	return true
}

// Line returns a treeLine to be used to represent the node on the screen.
func (node *TreeNode) Line() treeLine {
	indent := strings.Repeat(nodeIndent, node.Depth()-1)

	color := tcell.ColorReset
	if node.selected {
		color = tcell.GetColor("blue")
	}

	return treeLine{
		text:  fmt.Sprintf("%s%s", indent, node.task.Title),
		color: color,
	}
}

// newTreeNode returns a TreeNode with the given parent and Task.
func newTreeNode(parent *TreeNode, task *Task) *TreeNode {
	return &TreeNode{
		task:     task,
		parent:   parent,
		children: []*TreeNode{},
	}
}

func rootTreeNode() *TreeNode {
	return &TreeNode{
		task:     nil,
		parent:   nil,
		children: []*TreeNode{},
	}
}

// TreeLine is a line in the screen representation of a tree.
type treeLine struct {
	text  string
	color tcell.Color
}

// Tree is the tree of tasks in Impulse.
//
// It's largely adapted from tview.TreeView.
type Tree struct {
	// pull in default functions to satisfy the tview.Primitive interface
	*tview.Box

	root *TreeNode
}

// PushTask adds to the root node a new TreeNode for the given task.
//
// The new TreeNode, which is returned, will be the topmost child of tree's root node.
func (tree *Tree) PushTask(task *Task) *TreeNode {
	return tree.root.PushTask(task)
}

// Walk calls walk on the tree's root node.
func (tree *Tree) Walk(callback func(*TreeNode, *TreeNode) bool) bool {
	return tree.root.Walk(callback)
}

// Draw draws this primitive onto the screen.
func (t *Tree) Draw(screen tcell.Screen) {
	lines := []treeLine{}

	t.Walk(func(node, parent *TreeNode) bool {
		if parent == nil {
			// Root node has no task to draw
			return true
		}
		lines = append(lines, node.Line())
		return true
	})

	_, _, width, _ := t.GetInnerRect()
	nLines := len(lines)
	for i, line := range lines {
		tview.Print(
			screen,
			line.text,
			0,               // x
			nLines-1-i,      // y
			width,           // maxWidth
			tview.AlignLeft, // align
			line.color,      // color
		)
	}
}

// this is the one we have to implement
func (tree *Tree) InputHandler() func(*tcell.EventKey, func(p tview.Primitive)) {
	return func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {}
}

// NewTree returns a new Tree with a root node.
func NewTree() *Tree {
	return &Tree{
		Box:  tview.NewBox(),
		root: rootTreeNode(),
	}
}
