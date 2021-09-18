package common

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTreeNode_AddChild(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	n := NewTreeNode("")
	n.AddChild(NewTreeNode("a"))
	n.AddChild(NewTreeNode("b"))
	n.AddChild(NewTreeNode("c"))

	rslt := make([]string, 0)
	for _, ch := range n.Children {
		rslt = append(rslt, ch.Referent)
	}
	assert.Equal([]string{"a", "b", "c"}, rslt)
}

func TestTreeNode_InsertChild(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	n := NewTreeNode("")
	n.AddChild(NewTreeNode("b"))
	n.AddChild(NewTreeNode("d"))
	n.InsertChild(0, NewTreeNode("a"))
	n.InsertChild(2, NewTreeNode("c"))
	n.InsertChild(4, NewTreeNode("e"))

	rslt := make([]string, 0)
	for _, ch := range n.Children {
		rslt = append(rslt, ch.Referent)
	}
	assert.Equal([]string{"a", "b", "c", "d", "e"}, rslt)
}

func TestTreeNode_InsertChild_Empty(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	n := NewTreeNode("")
	n.InsertChild(0, NewTreeNode("a"))

	rslt := make([]string, 0)
	for _, ch := range n.Children {
		rslt = append(rslt, ch.Referent)
	}
	assert.Equal([]string{"a"}, rslt)
}

func TestTreeNode_Walk(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	mp := MakePasta()[0]
	referents := make([]string, 0)
	err := mp.RootNode.Walk(func(vn *TreeNode) error {
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

func TestTreeNode_Walk_SkipSubtree(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	mp := MakePasta()[0]
	referents := make([]string, 0)
	err := mp.RootNode.Walk(func(vn *TreeNode) error {
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

func TestTreeNode_Walk_Error(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	mp := MakePasta()[0]
	referents := make([]string, 0)
	expErr := errors.New("bleuhgarrg")
	err := mp.RootNode.Walk(func(vn *TreeNode) error {
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

func TestTreeNode_Walk_ZeroChildren(t *testing.T) {
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

func TestTreeNode_Equal(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	a := MakePasta()[0]
	b := MakePasta()[0]
	assert.True(a.Equal(b))
	assert.True(b.Equal(a))
}

// Two trees should not be evaluated as Equal if one of the corresponding node pairs differs in
// referent
func TestTreeNode_Equal_Not_Referent(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	a := MakePasta()[0]
	b := MakePasta()[0]
	b.RootNode.Walk(func(n *TreeNode) error {
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
func TestTreeNode_Equal_Not_ChildCount(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	a := MakePasta()[0]
	b := MakePasta()[0]
	b.RootNode.Walk(func(n *TreeNode) error {
		if n.Referent == "put pot on burner" {
			n.AddChild(NewTreeNode("fhgwhgads"))
		}
		return nil
	})
	assert.False(a.Equal(b))
	assert.False(b.Equal(a))
}

func TestTreeNode_Equal_ZeroChildren(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	a := NewTreeNode("childless node")
	b := NewTreeNode("childless node")
	assert.True(a.Equal(b))
	assert.True(b.Equal(a))
}

func TestTreeNode_Depth(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	a := NewTreeNode("grandparent")
	b := NewTreeNode("parent")
	c := NewTreeNode("child")
	b.AddChild(c)
	a.AddChild(b)

	assert.Equal(0, a.Depth())
	assert.Equal(1, b.Depth())
	assert.Equal(2, c.Depth())
}
