package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danslimmon/impulse/common"
)

type GetTaskListResponse struct {
	Error  string         `json:"error,omitempty"`
	Result []*common.Task `json:"result,omitempty"`
}

func (apiResp *GetTaskListResponse) UnmarshalJSON(b []byte) error {
	// Declare a local type so that we don't get infinite recursion when we call json.Unmarshal
	type gtlr GetTaskListResponse
	transient := new(gtlr)
	err := json.Unmarshal(b, transient)
	if err != nil {
		return err
	}
	*apiResp = GetTaskListResponse(*transient)

	// Parent field is not included in JSON representation of *comon.TreeNode, so we must walk each
	// Task's tree and populate Parent.
	for _, task := range apiResp.Result {
		err := task.RootNode.Walk(func(tn *common.TreeNode) error {
			for _, child := range tn.Children {
				child.Parent = tn
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

type ImpulseAPI struct {
	server    *http.Server
	taskstore Taskstore
}

// assignTaskstore obtains a default Taskstore implementation if one is not already injected.
func (api *ImpulseAPI) assignTaskstore() {
	if api.taskstore != nil {
		// We already have a Taskstore implementation dependency-injected.
		return
	}

	ds := NewFilesystemDatastore(DataDir)
	ts := NewBasicTaskstore(ds)
	api.taskstore = ts
}

func (api *ImpulseAPI) listenAndServe(addr string, router http.Handler) {
	api.server.ListenAndServe()
}

// Start starts the Impulse API server, which will listen for requests until Stop is called.
func (api *ImpulseAPI) Start(addr string) error {
	api.assignTaskstore()

	router := gin.Default()
	router.GET("/tasklist/:tasklistname", func(c *gin.Context) {
		name := c.Param("tasklistname")
		tasks, err := api.taskstore.GetList(name)
		if err != nil {
			//@DEBUG
			fmt.Printf("x-bravo error: %s\n", err.Error())
			c.JSON(404, GetTaskListResponse{Error: err.Error()})
		} else {
			//@DEBUG
			fmt.Printf("x-bravo no error\n")
			c.JSON(200, GetTaskListResponse{Result: tasks})
		}
	})

	api.server = &http.Server{
		Addr:    addr,
		Handler: router,
	}
	// we do Listen and Serve separately to make sure that the listener is open before this function
	// returns. otherwise, the tests will sometimes get Connection Refused because of the race
	// condition.
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	go api.server.Serve(listener)
	return nil
}

// Stop stops the Impulse API server.
func (api *ImpulseAPI) Stop() error {
	return api.server.Shutdown(context.Background())
}

func NewImpulseAPI(ts Taskstore) *ImpulseAPI {
	return &ImpulseAPI{
		taskstore: ts,
	}
}
