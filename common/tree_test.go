package common

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Returns a "make pasta" task.
//
// Should be identical to the contents of server/testdata/make_pasta
//
//             put water in pot
//             put pot on burner
//             turn burner on
//         boil water
//         put pasta in water
//         [b cooked]
//         drain pasta
//     make pasta
func MakePasta() *TreeNode {
	makePasta := NewTreeNode("make pasta")

	boilWater := NewTreeNode("boil water")
	boilWater.AddChild(NewTreeNode("put water in pot"))
	boilWater.AddChild(NewTreeNode("put pot on burner"))
	boilWater.AddChild(NewTreeNode("turn burner on"))
	makePasta.AddChild(boilWater)

	makePasta.AddChild(NewTreeNode("put pasta in water"))
	makePasta.AddChild(NewTreeNode("[b cooked]"))
	makePasta.AddChild(NewTreeNode("drain pasta"))

	return makePasta
}

func TestTree_Walk(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	n := MakePasta()
	referents := make([]string, 0)
	err := n.Walk(func(vn *TreeNode) error {
		referents = append(referents, vn.Referent)
		return nil
	})

	assert.Nil(err)
	assert.Equal(
		[]string{
			"make pasta",
			"boil water",
			"put water in pot",
			"put pot on burner",
			"turn burner on",
			"put pasta in water",
			"[b cooked]",
			"drain pasta",
		},
		referents,
	)
}

func TestTree_Walk_SkipSubtree(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	n := MakePasta()
	referents := make([]string, 0)
	err := n.Walk(func(vn *TreeNode) error {
		referents = append(referents, vn.Referent)
		if vn.Referent == "boil water" {
			return SkipSubtree
		}
		return nil
	})

	assert.Nil(err)
	assert.Equal(
		[]string{
			"make pasta",
			"boil water",
			"put pasta in water",
			"[b cooked]",
			"drain pasta",
		},
		referents,
	)
}

func TestTree_Walk_Error(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	n := MakePasta()
	referents := make([]string, 0)
	expErr := errors.New("bleuhgarrg")
	err := n.Walk(func(vn *TreeNode) error {
		referents = append(referents, vn.Referent)
		if vn.Referent == "boil water" {
			return expErr
		}
		return nil
	})

	assert.Equal(expErr, err)
	assert.Equal(
		[]string{
			"make pasta",
			"boil water",
		},
		referents,
	)
}

func TestTree_Walk_ZeroChildren(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	n := NewTreeNode("childless node")
	referents := make([]string, 0)
	err := n.Walk(func(vn *TreeNode) error {
		referents = append(referents, vn.Referent)
		return nil
	})

	assert.Equal(nil, err)
	assert.Equal(
		[]string{
			"childless node",
		},
		referents,
	)
}

func TestTree_Equal(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	a := MakePasta()
	b := MakePasta()
	assert.True(a.Equal(b))
	assert.True(b.Equal(a))
}

// Two trees should not be evaluated as Equal if one of the corresponding node pairs differs in
// referent
func TestTree_Equal_Not_Referent(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	a := MakePasta()
	b := MakePasta()
	b.Walk(func(n *TreeNode) error {
		if n.Referent == "put pot on burner" {
			n.Referent = "put pot in bong"
		}
		return nil
	})
	assert.False(a.Equal(b))
	assert.False(b.Equal(a))
}

// Two trees should not be evaluated as Equal if one of the corresponding node pairs differs in
// the number of children they have
func TestTree_Equal_Not_ChildCount(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	a := MakePasta()
	b := MakePasta()
	b.Walk(func(n *TreeNode) error {
		if n.Referent == "put pot on burner" {
			n.AddChild(NewTreeNode("fhgwhgads"))
		}
		return nil
	})
	assert.False(a.Equal(b))
	assert.False(b.Equal(a))
}

func TestTree_Equal_ZeroChildren(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	a := NewTreeNode("childless node")
	b := NewTreeNode("childless node")
	assert.True(a.Equal(b))
	assert.True(b.Equal(a))
}
