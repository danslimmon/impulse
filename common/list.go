package common

import (
	"bytes"
)

type BlopList struct {
	Name  string  `json:"name"`
	Blops []*Blop `json:"blops"`
}

// blipLine is a line from a marshaled blop
type blipLine struct {
	// Number of tabs the line is indented by
	Indent int
	// Text after indent
	Text string
}

// textToBlipLine converts the given line from a marshaled Blop into a blipLine struct
func textToBlipLine(text []byte) (blipLine, error) {
	rslt := blipLine{}
	for i := range text {
		if string(text[i]) == "\t" {
			rslt.Indent = rslt.Indent + 1
		} else {
			rslt.Text = string(text[i:])
			break
		}
	}
	return rslt, nil
}

// UnmarshalBlopList reads the given byte array and returns a corresponding *BlopList named name.
func UnmarshalBlopList(name string, b []byte) (*TreeNode, error) {
	rslt := new(BlopList)
	rslt.Name = name

	// Convert the text lines into blipLine structs
	lines := bytes.Split(b, []byte("\n"))
	blipLines := make([]blipLine, len(lines))
	for i := range lines {
		blipLine, err := textToBlipLine(lines[i])
		if err != nil {
			return nil, err
		}
		blipLines[i] = blipLine
	}

	// Split blipLines into groups, where each group corresponds to a blip. For example, if
	// blipLines is derived from
	//
	//     reply to slack message
	//         feed dog
	//             open can
	//         feed cat
	//     feed pets
	//
	// then at the end of this stanza, groups will be equal to
	//
	//     [][]blipLine{
	//         []blipLine{
	//             blipLine{0, "feed pets"},
	//             blipLine{1, "feed cat"},
	//             blipLine{2, "open can"},
	//             blipLine{1, "feed dog"},
	//         },
	//		   []blipLine{
	//             blipLine{0, "reply to slack message"},
	//         },
	//     }
	breakpoints := make([]int, 0)
	for i := len(blipLines) - 1; i >= 0; i-- {
		curLine := blipLines[i]
		if curLine.Indent == 0 {
			breakpoints = append(breakpoints, i)
		}
	}
	groups := make([][]blipLine, len(breakpoints))
	for bpi := 0; bpi < len(breakpoints) - 1; bpi++ {
		groups[bpi] = blipLines[breakpoints[bpi]:breakpoints[bpi+1]]
	}
	groups = append(groups, blipLines[breakpoints[len(breakpoints)-1]:]

	// Convert blipLine groups into Blops
	for _, g := range groups {
		blop := new(Blop)
		for _, line := range g {
			blop.Text = g
			if line.Indentuu
		}
	}

	return blopList, nil
}
