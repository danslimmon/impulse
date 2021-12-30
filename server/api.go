package server

import (
	"context"
	"encoding/json"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/danslimmon/impulse/common"
)

type GetTaskListResponse struct {
	Error  string         `json:"error,omitempty"`
	Result []*common.Task `json:"result,omitempty"`
}

type ArchiveLineResponse struct {
	Error string `json:"error,omitempty"`
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

type Server struct {
	server    *http.Server
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

func (api *Server) listenAndServe(addr string, router http.Handler) {
	api.server.ListenAndServe()
}

// Start starts the Impulse API server, which will listen for requests until Stop is called.
func (api *Server) Start(addr string) error {
	api.assignTaskstore()

	router := gin.Default()

	router.GET("/tasklist/:tasklistname", func(c *gin.Context) {
		name := c.Param("tasklistname")
		tasks, err := api.taskstore.GetList(name)
		if err != nil {
			c.JSON(404, GetTaskListResponse{Error: err.Error()})
		} else {
			c.JSON(200, GetTaskListResponse{Result: tasks})
		}
	})

	router.POST("/archive_line/:line_id", func(c *gin.Context) {
		lineId := common.LineID(c.Param("line_id"))
		err := api.taskstore.ArchiveLine(lineId)
		if err != nil {
			c.JSON(404, ArchiveLineResponse{Error: err.Error()})
		} else {
			c.JSON(200, ArchiveLineResponse{})
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
func (api *Server) Stop() error {
	return api.server.Shutdown(context.Background())
}

func NewServer(ts Taskstore) *Server {
	return &Server{
		taskstore: ts,
	}
}
