package server

import (
	"github.com/danslimmon/impulse/common"
)

// Unmarshal reads the given data and constructs a corresponding *TreeNode.
//
// The *TreeNode we return has a nil referent and resides one level above the bottommost line of the
// data. The bottommost line and any siblings of it are assumed to be Blops. For example, if the
// data we read is
//
//     check twitter
//         boil water
//     make pasta
//
// then the returned *TreeNode's children will refer to &Blop{"check twitter"} and &Blop{"make
// pasta"}, in that order.
func Unmarshal(data []byte) (*common.TreeNode, error) {
	return nil, nil
}

// unmarshalSection reads the given data constructs a corresponding *TreeNode.
//
// The *TreeNode we return has a nil referent and resides one level above the bottommost line of the
// data. The referent of this *TreeNode, and of all its descendants, is a *Blip.
func unmarshalSection(data []byte) (*common.TreeNode, error) {
	return nil, nil
}
