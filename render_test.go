package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Render returns a string representation of the root Frame
func TestRender(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := NewFrame("")
	c := NewFrame("cook pasta")
	f.Push(c)
	gc := NewFrame("boil water")
	c.Push(gc)
	gc.Insert(NewFrame("put water in pot"))
	gc.Insert(NewFrame("put pot on burner"))
	gc.Insert(NewFrame("turn on burner"))
	f.Push(NewFrame("check twitter"))

	exp := `川 check twitter
        / put water in pot
        / put pot on burner
        / turn on burner
    / boil water
口 cook pasta`
	assert.Equal(exp, Render(f))
}

// Render panics if passed a non-root Frame.
func TestRender_NotRootFrame(t *testing.T) {
	t.Fail()
}

// renderLine returns a single-line string representation of a given Frame.
func TestRenderLine(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := NewFrame("")
	f.Push(NewFrame("cook pasta"))
	f.Children[0].Push(NewFrame("boil water"))
	f.Children[0].Children[0].Push(NewFrame("put water in pot"))
	f.Push(NewFrame("check twitter"))

	assert.Equal("川 check twitter", renderLine(f.Children[1]))
	assert.Equal("        / put water in pot", renderLine(f.Children[0].Children[0].Children[0]))
	assert.Equal("    / boil water", renderLine(f.Children[0].Children[0]))
	assert.Equal("口 cook pasta", renderLine(f.Children[0]))
}

// renderLine panics when passed the root Frame.
func TestRenderLine_RootFrame(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := NewFrame("")
	assert.Panics(func() {
		renderLine(f)
	})
}
