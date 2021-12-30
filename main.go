package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/danslimmon/impulse/client"
	"github.com/danslimmon/impulse/common"
	"github.com/danslimmon/impulse/server"
)

func main() {
	dataDir := os.Getenv("IMPULSE_DATADIR")
	if dataDir == "" {
		panic("error: IMPULSE_DATADIR environment variable required")
	}

	addr := "127.0.0.1:30271"
	ds := server.NewFilesystemDatastore(dataDir)
	ts := server.NewBasicTaskstore(ds)
	apiServer := server.NewServer(ts)
	if err := apiServer.Start(addr); err != nil {
		panic("failed to start test server on " + addr + ": " + err.Error())
	}

	apiClient := client.NewClient(addr)
	switch os.Args[1] {
	case "show":
		resp, err := apiClient.GetTaskList(os.Args[2])
		if err != nil {
			panic(fmt.Sprintf("failed to get task list `%s`: %s", os.Args[2], err.Error()))
		}

		for _, t := range resp.Result {
			t.RootNode.WalkFromTop(func(n *common.TreeNode) error {
				fmt.Printf(
					"%s%v\n",
					strings.Repeat("    ", n.Depth()),
					n.Referent,
				)
				return nil
			})
		}
	case "archive":
		resp, err := apiClient.ArchiveLine(os.Args[2], os.Args[3])
	}
}
