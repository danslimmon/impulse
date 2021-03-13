package main

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
