package server

import (
	"github.com/danslimmon/impulse/common"
)

type GetTaskRequest struct {
	TaskIDs []string
}

type GetTaskResponse struct {
	Response
	Tasks []*common.Task
}

// GetTask retrieves the tasks whose task IDs are those specified by req.TaskIDs. The response's
// Tasks attribute contains the tasks returned, in the same order as they were listed in TaskIDs.
//
// Should any of the IDs in TaskIDs not match a task known by the server, it's an error. Duplicate
// TaskIDs values are okay – resp.Tasks will just contain the same task however many times its ID is
// repeated in the input.
func (s *Server) GetTask(req *GetTaskRequest, resp *GetTaskResponse) error {
	resp.Tasks = make([]*common.Task, len(req.TaskIDs))
	for i := range req.TaskIDs {
		task, err := s.taskstore.GetTask(common.LineID(req.TaskIDs[i]))
		if err != nil {
			return err
		}
		resp.Tasks[i] = task
	}

	return nil
}
