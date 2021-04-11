package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// makePastaTree returns a tree of tasks for making pasta:
//
//         drain pasta
//         wait for pasta to cook
//             put water in pot
//             put pot on burner
//             turn on burner
//         boil water
//     make pasta
// root
func makePastaTree() *Tree {
	tree := NewTree()
	makePasta := tree.PushTask(&Task{Title: "make pasta"})
	boilWater := makePasta.PushTask(&Task{Title: "boil water"})
	boilWater.PushTask(&Task{Title: "turn on burner"})
	boilWater.PushTask(&Task{Title: "put pot on burner"})
	boilWater.PushTask(&Task{Title: "put water in pot"})
	makePasta.PushTask(&Task{Title: "wait for pasta to cook"})
	makePasta.PushTask(&Task{Title: "drain pasta"})
	return tree
}

func TestTreeNode_PushTask(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := rootTreeNode()
	root.PushTask(&Task{Title: "make pasta"})
	assert.Equal(1, len(root.children))
	assert.Equal("make pasta", root.children[0].task.Title)
	assert.Equal(root, root.children[0].parent)

	makePasta := root.children[0]
	makePasta.PushTask(&Task{Title: "boil water"})
	assert.Equal(1, len(makePasta.children))
	assert.Equal("boil water", makePasta.children[0].task.Title)
	assert.Equal(makePasta, makePasta.children[0].parent)

	boilWater := makePasta.children[0]
	boilWater.PushTask(&Task{Title: "turn on burner"})
	boilWater.PushTask(&Task{Title: "put pot on burner"})
	boilWater.PushTask(&Task{Title: "put water in pot"})
	assert.Equal(3, len(boilWater.children))
	assert.Equal("put water in pot", boilWater.children[2].task.Title)
	assert.Equal("put pot on burner", boilWater.children[1].task.Title)
	assert.Equal("turn on burner", boilWater.children[0].task.Title)
	for _, child := range boilWater.children {
		assert.Equal(boilWater, child.parent)
	}
}

func TestTreeNode_Walk(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	type invocation struct {
		Node   *TreeNode
		Parent *TreeNode
	}
	invocations := make([]invocation, 0)
	record := func(node, parent *TreeNode) bool {
		invocations = append(invocations, invocation{node, parent})
		return true
	}

	tree := makePastaTree()
	tree.Walk(record)

	assert.Equal(8, len(invocations))
	assert.Equal(invocations[0].Node, tree.root)
	assert.Equal(invocations[7].Node.task.Title, "drain pasta")
}

// When callback returns false, Walk should not recurse further.
func TestTreeNode_Walk_Stop(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	type invocation struct {
		Node   *TreeNode
		Parent *TreeNode
	}
	invocations := make([]invocation, 0)
	cb := func(node, parent *TreeNode) bool {
		invocations = append(invocations, invocation{node, parent})
		if node.task != nil && node.task.Title == "boil water" {
			// Do not recurse into the "boil water" task
			return false
		}
		return true
	}

	tree := makePastaTree()
	tree.Walk(cb)

	assert.Equal(5, len(invocations))
	assert.Equal(invocations[0].Node, tree.root)
	assert.Equal(invocations[4].Node.task.Title, "drain pasta")
}
