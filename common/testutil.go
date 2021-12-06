package common

// Returns a "make pasta" tasklist.
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
func MakePasta() []*Task {
	makePasta := NewTreeNode("make pasta")

	boilWater := NewTreeNode("boil water")
	boilWater.AddChild(NewTreeNode("put water in pot"))
	boilWater.AddChild(NewTreeNode("put pot on burner"))
	boilWater.AddChild(NewTreeNode("turn burner on"))
	makePasta.AddChild(boilWater)

	makePasta.AddChild(NewTreeNode("put pasta in water"))
	makePasta.AddChild(NewTreeNode("[b cooked]"))
	makePasta.AddChild(NewTreeNode("drain pasta"))

	return []*Task{NewTask(makePasta)}
}

// Returns a "multiple nested" tasklist.
//
// Should be identical to the contents of server/testdata/multiple_nested
//
//         subsubtask 0
//     subtask 0
// task 0
//     subtask 1
// task 1
func MultipleNested() []*Task {
	task0 := NewTreeNode("task 0")
	subtask0 := NewTreeNode("subtask 0")
	subsubtask0 := NewTreeNode("subsubtask 0")
	subtask0.AddChild(subsubtask0)
	task0.AddChild(subtask0)

	task1 := NewTreeNode("task 1")
	subtask1 := NewTreeNode("subtask 1")
	task1.AddChild(subtask1)

	return []*Task{
		NewTask(task0),
		NewTask(task1),
	}
}
