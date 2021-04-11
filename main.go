package main

import (
	"github.com/rivo/tview"
)

// defaultTree returns the tree we're using as sample data at this early moment in development
func defaultTree() *Tree {
	tree := NewTree()
	makePasta := tree.PushTask(&Task{Title: "make pasta"})
	makePasta.PushTask(&Task{Title: "boil water"})
	return tree
}

func main() {
	tree := defaultTree()
	if err := tview.NewApplication().SetRoot(tree, true).Run(); err != nil {
		panic(err)
	}
}
