package main

import (
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
	assert.Equal("cook pasta", c.Name)
	assert.Equal(f, c.Parent)

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
}
