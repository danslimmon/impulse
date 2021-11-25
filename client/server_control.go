package main

import (
	"fmt"

	"github.com/danslimmon/impulse/server"
)

var serverHandle *server.ImpulseAPI

// StartServer starts an impulse server listening on the given IP:port pair.
//
// dataDir is the path to where the server's data lives
func StartServer(addr, dataDir string) error {
	if serverHandle != nil {
		return fmt.Errorf("server already started")
	}

	// got a plan to actually split apart the server and client… until then, this nonsense
	ds := server.NewFilesystemDatastore("_")
	ts := server.NewBasicTaskstore(ds)
	serverHandle = server.NewImpulseAPI(ts)
	return serverHandle.Start(addr)
}

// StopServer stops the impulse server previously started with StartServer.
func StopServer() error {
	if serverHandle == nil {
		return fmt.Errorf("server already stopped")
	}

	return serverHandle.Stop()
}
