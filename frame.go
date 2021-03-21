package main

import (
	"errors"
)

// FrameWalkFunc is a function that can be passed to Frame.Walk.
//
// At each Frame visited by Walk, the FrameWalkFunc is passed that Frame. If the FrameWalkFunc
// returns an error, the walk will be aborted and that error will be returned by Walk.
//
// FrameWalkFunc should surface any error it receives in the course of its operation by returning
// it, so that Walk can return that to its caller. If FrameWalkFunc wants to end the walk early but
// there's not an error per se, then it should return the special error EndWalk. If the caller of
// Walk passes a FrameWalkFunc that may return EndWalk, then the caller may check whether Walk
// returned EndWalk and act accordingly.
type FrameWalkFunc func(*Frame) error

var (
	EndWalk error = errors.New("Ended walk at FrameWalkFunc's request")
)

// Frame is a node of the stack. It represents a task.
//
// You can also just think of a Frame as a single line in the stack view.
//
// A Frame has a parent Frame, unless it is the root Frame (described below).
//
// A Frame has a nonnegative integer number of ordered child Frames. Child Frames are ordered from
// low to high in the stack (i.e. the last child Frame is the uppermost in the stack). When a
// Frame's Pop method is called, its uppermost child Frame is removed. When a Frame's Push method is
// called, a new child Frame is added at the uppermost position.
//
// A special case is the root Frame. The root Frame's parent is nil. When the root Frame's last
// child is popped, impulse exits.
type Frame struct {
	Name string

	Parent   *Frame
	Children []*Frame
}

// Push adds the Frame c to the Frame f's children.
//
// When Push returns, c is at the uppermost position among f's children.
func (f *Frame) Push(c *Frame) {
	c.Parent = f
	f.Children = append(f.Children, c)
}

// Pop removes f's uppermost child and returns it.
//
// If f has no children, Pop returns nil.
func (f *Frame) Pop() *Frame {
	if len(f.Children) == 0 {
		return nil
	}
	c := f.Children[len(f.Children)-1]
	f.Children = f.Children[:len(f.Children)-1]
	return c
}

// Insert adds c to f's children at the bottom-most position.
func (f *Frame) Insert(c *Frame) {
	f.Children = append([]*Frame{c}, f.Children...)
}

// Depth returns the number of ancestors Frame has.
//
// The root frame's Depth is 0. Children of the root frame have Depth 1. And so on.
func (f *Frame) Depth() int {
	if f.Parent == nil {
		return 0
	}
	return f.Parent.Depth() + 1
}

// Walk calls the given function once for each descendant frame of f.
//
// The order is depth-first, bottom to top.
//
// If fn returns an error, the walk is aborted and that error is returned by Walk. See
// FrameWalkFunc's definition for more info.
func (f *Frame) Walk(fn FrameWalkFunc) error {
	err := fn(f)
	if err != nil {
		return err
	}

	for _, c := range f.Children {
		err := c.Walk(fn)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewFrame returns a new Frame with the given name.
//
// The Frame returned has nil for a parent. If it is pushed onto a parent Frame, Push will assign a
// pointer to the parent Frame to the pushed Frame's Parent field.
func NewFrame(name string) *Frame {
	return &Frame{
		Name:     name,
		Parent:   nil,
		Children: []*Frame{},
	}
}
