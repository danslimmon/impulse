package common

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
