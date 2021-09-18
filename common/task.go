package common

type Task struct {
	RootNode *TreeNode
}

// Equal determines whether the tasks a and b are equal.
func (a *Task) Equal(b *Task) bool {
	return a.RootNode.Equal(b.RootNode)
}

func NewTask(n *TreeNode) *Task {
	return &Task{RootNode: n}
}
