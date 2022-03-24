package server

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
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

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
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
