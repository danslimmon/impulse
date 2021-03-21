package main

import (
	"fmt"
	"strings"
)

// Render returns a string representation of the root Frame.
//
// If f is not the root Frame, Render panics.
func Render(f *Frame) string {
	// At the end of this stanza, frames is a list of frames with frames[0] being the frame that
	// will be displayed at the bottom and frames[-1] being the frame that will be displayed at the
	// top.
	frames := make([]*Frame, 0)
	err := f.Walk(func(node *Frame) error {
		if node.Depth() == 0 {
			// Exclude the root Frame from frames, since the root Frame is not displayed as a line.
			return nil
		}
		frames = append(frames, node)
		return nil
	})
	if err != nil {
		panic("Unexpected error encountered during Walk")
	}

	// Reverse the order of frames so that the topmost frame is frames[-1] and the bottom one is
	// frames[0]. This is the order in which we want to output the frames' string representations.
	for i, j := 0, len(frames)-1; i < j; i, j = i+1, j-1 {
		frames[i], frames[j] = frames[j], frames[i]
	}

	lines := make([]string, len(frames))
	for i, f := range frames {
		lines[i] = renderLine(f)
	}

	return strings.Join(lines, "\n")
}

// renderLine returns a string representation of the given Frame.
//
// The string representation of a Frame has the following anatomy:
//
//     [INDENT][ICON] [NAME]
//
// INDENT is a number of space characters equal to 4 times (f.Depth() - 1). Children of the root
// Frame have INDENT = "", grandchildren of the root Frame have INDENT = "    ", etc.
//
// ICON is equal to "川" if f is the topmost frame of depth 1. ICON equals "口" if f has depth 1 but
// is not the topmost frame with depth 1. In all other cases, ICON equals "/".
//
// NAME is equal to f.Name.
//
// The root Frame is not represented to the user visually. As such renderLine when called on the
// root Frame will panic.
func renderLine(f *Frame) string {
	if f.Depth() == 0 {
		panic("renderLine called on root Frame")
	}

	indent := strings.Repeat("    ", f.Depth()-1)

	icon := "/"
	if f.Depth() == 1 {
		if f == f.Parent.Children[len(f.Parent.Children)-1] {
			icon = "川"
		} else {
			icon = "口"
		}
	}

	return fmt.Sprintf(
		"%s%s %s",
		indent,
		icon,
		f.Name,
	)
}
