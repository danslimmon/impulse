package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrame_Push(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := NewFrame("")
	assert.Nil(f.Parent)

	f.Push(NewFrame("cook pasta"))
	assert.Equal(1, len(f.Children))
	c := f.Children[0]
	assert.Equal("cook pasta", c.Name)
	assert.Equal(f, c.Parent)

	c.Push(NewFrame("turn on burner"))
	c.Push(NewFrame("put pot on burner"))
	c.Push(NewFrame("put water in pot"))
	assert.Equal(3, len(c.Children))
	assert.Equal("put water in pot", c.Children[2].Name)
	assert.Equal("put pot on burner", c.Children[1].Name)
	assert.Equal("turn on burner", c.Children[0].Name)
	for _, gc := range c.Children {
		assert.Equal(c, gc.Parent)
		assert.Equal(f, gc.Parent.Parent)
		assert.Nil(gc.Parent.Parent.Parent)
	}
}

func TestFrame_Pop(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := NewFrame("")
	assert.Nil(f.Parent)
	c := NewFrame("cook pasta")
	f.Push(c)

	c.Push(NewFrame("turn on burner"))
	c.Push(NewFrame("put pot on burner"))
	c.Push(NewFrame("put water in pot"))

	gc := c.Pop()
	assert.NotNil(gc)
	assert.Equal("put water in pot", gc.Name)
	assert.Equal(c, gc.Parent)

	gc = c.Pop()
	assert.NotNil(gc)
	assert.Equal("put pot on burner", gc.Name)
	assert.Equal(c, gc.Parent)

	gc = c.Pop()
	assert.NotNil(gc)
	assert.Equal("turn on burner", gc.Name)
	assert.Equal(c, gc.Parent)

	assert.Nil(c.Pop())
	assert.Equal(c, f.Pop())
	assert.Nil(f.Pop())
	assert.Nil(f.Pop())
}

func TestFrame_Insert(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := NewFrame("cook pasta")
	f.Insert(NewFrame("put water in pot"))
	f.Insert(NewFrame("put pot on burner"))
	f.Insert(NewFrame("turn on burner"))

	c := f.Pop()
	assert.NotNil(c)
	assert.Equal(c.Name, "put water in pot")

	c = f.Pop()
	assert.NotNil(c)
	assert.Equal(c.Name, "put pot on burner")

	c = f.Pop()
	assert.NotNil(c)
	assert.Equal(c.Name, "turn on burner")
}

func TestFrame_Walk(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := NewFrame("")
	f.Push(NewFrame("cook pasta"))
	f.Children[0].Push(NewFrame("boil water"))

	assert.Equal(0, f.Depth())
	assert.Equal(1, f.Children[0].Depth())
	assert.Equal(2, f.Children[0].Children[0].Depth())
}

func TestFrame_Depth(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := NewFrame("")
	f.Push(NewFrame("cook pasta"))
	f.Children[0].Push(NewFrame("boil water"))
	f.Children[0].Children[0].Push(NewFrame("turn on burner"))
	f.Children[0].Children[0].Push(NewFrame("put pot on burner"))
	f.Children[0].Children[0].Push(NewFrame("put water in pot"))
	f.Push(NewFrame("check twitter"))

	exp := []string{
		"",
		"cook pasta",
		"boil water",
		"turn on burner",
		"put pot on burner",
		"put water in pot",
		"check twitter",
	}
	var i int
	err := f.Walk(func(node *Frame) error {
		assert.Equal(exp[i], node.Name)
		i++
		return nil
	})
	assert.Equal(len(exp), i)
	assert.Nil(err)
}

// Walk should exit immediately and return the error that FrameWalkFunc returned.
func TestFrame_Walk_Error(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := NewFrame("cook pasta")
	f.Push(NewFrame("boil water"))

	var i int
	errToReturn := fmt.Errorf("dumb error")
	err := f.Walk(func(node *Frame) error {
		i++
		return errToReturn
	})
	assert.Equal(errToReturn, err)
	assert.Equal(1, i)
}
