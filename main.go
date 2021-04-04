package main

import (
	"github.com/rivo/tview"
)

// defaultTree returns the tree we're using as sample data at this early moment in development
func defaultTree() *Tree {
	tree := NewTree()

	var makePasta, boilWater *TreeNode
	tree.root.children = []*TreeNode{makePasta}
	makePasta = &TreeNode{
		task:     &Task{Title: "make pasta"},
		parent:   tree.root,
		children: []*TreeNode{boilWater},
	}
	boilWater = &TreeNode{
		task:     &Task{Title: "boil water"},
		parent:   makePasta,
		children: []*TreeNode{},
	}

	return tree
}

func main() {
	tree := defaultTree()
	if err := tview.NewApplication().SetRoot(tree, true).Run(); err != nil {
		panic(err)
	}
}
