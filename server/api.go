package server

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/danslimmon/impulse/common"
)

type Server struct {
	taskstore Taskstore
}

// assignTaskstore obtains a default Taskstore implementation if one is not already injected.
func (api *Server) assignTaskstore() {
	if api.taskstore != nil {
		// We already have a Taskstore implementation dependency-injected.
		return
	}

	ds := NewFilesystemDatastore(DataDir)
	ts := NewBasicTaskstore(ds)
	api.taskstore = ts
}

// Start starts the Impulse API server, which will listen for requests until Stop is called.
func (api *Server) Start(addr string) error {
	api.assignTaskstore()
	rpc.Register(api)

	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:4358")
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			jsonrpc.ServeConn(conn)
		}
	}()

	return nil
}

func NewServer(ts Taskstore) *Server {
	return &Server{
		taskstore: ts,
	}
}

type Response struct {
	StateID string
}

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
// Should any of the IDs in TaskIDs be unknown by the server, it's an error. Duplicate TaskIDs
// values are okay – resp.Tasks will just contain the same task however many times its ID is repeated
// in the input.
func (s *Server) GetTask(req *GetTaskRequest, resp *GetTaskResponse) error {
	/*
		if len(req.TaskIDs) == 0 {
			return nil
		}
	*/

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
