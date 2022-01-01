package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/danslimmon/impulse/common"
)

// ap is a singleton addrPool that we use to provision addrs for tests to listen on.
var ap *addrPool

// addrPool provisions unique addrs for tests to listen on.
//
// An addr is a string of the form "<IP>:<port>". We will provision up to 10000 unique addrs. If
// more than 10000 are requested, addrPool.Get panics.
type addrPool struct {
	port int
	mu   sync.Mutex
}

func (a *addrPool) Get() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.port == 0 {
		a.port = 30000
	} else if a.port >= 39999 {
		panic("more than 10000 addrs requested from addrPool")
	} else {
		a.port = a.port + 1
	}
	return fmt.Sprintf("127.0.0.1:%d", a.port)
}

// listenAddr returns a loopback address for a test instance of the API to listen on.
//
// The address is guaranteed not to be in use by any other test in the suite.
func listenAddr() string {
	if ap == nil {
		ap = new(addrPool)
	}
	return ap.Get()
}

func TestServer_GetTaskList(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ts := NewBasicTaskstore(
		NewFilesystemDatastore("./testdata"),
	)

	api := &Server{
		taskstore: ts,
	}
	addr := listenAddr()
	err := api.Start(addr)
	defer api.Stop()
	assert.Nil(err)

	resp, err := http.Get("http://" + addr + "/tasklist/make_pasta")
	assert.Nil(err)
	body, err := io.ReadAll(resp.Body)
	assert.Nil(err)

	apiResp := new(GetTaskListResponse)
	err = json.Unmarshal(body, apiResp)
	assert.Nil(err)
	assert.Equal("", apiResp.Error)
	assert.Equal(1, len(apiResp.Result))

	makePasta := common.MakePasta()[0]
	assert.True(apiResp.Result[0].RootNode.Equal(makePasta.RootNode))
}

func TestGetTaskListResponse_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	apiResp := new(GetTaskListResponse)
	b := []byte(`{
	"result": [
		{
			"tree": {
				"referent": "task title",
				"children": [
					{
						"referent": "first child"
					},
					{
						"referent": "second child"
					}
				]
			}
		}
	]
}`)

	err := json.Unmarshal(b, apiResp)
	assert.Nil(err)
	assert.Equal(1, len(apiResp.Result))
	task := apiResp.Result[0]

	assert.Nil(task.RootNode.Parent)
	for _, child := range task.RootNode.Children {
		assert.Equal(task.RootNode, child.Parent)
	}
}

func TestGetTaskListResponse_UnmarshalJSON_Error(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	apiResp := new(GetTaskListResponse)
	b := []byte(`{"error": "oh no there was an error"}`)

	err := json.Unmarshal(b, apiResp)
	assert.Nil(err)
	assert.Equal("oh no there was an error", apiResp.Error)
	assert.Nil(apiResp.Result)
}

func TestServer_ArchiveLine(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ds, cleanup := newFSDatastoreWithTestdata()
	defer cleanup()
	ts := NewBasicTaskstore(ds)

	api := &Server{
		taskstore: ts,
	}
	addr := listenAddr()
	err := api.Start(addr)
	defer api.Stop()
	assert.Nil(err)

	resp, err := http.Get(
		"http://" + addr + "/archive_line/" + url.PathEscape(string(
			common.GetLineID("make_pasta", "\t\tput water in pot"),
		)),
	)
	assert.Nil(err)
	body, err := io.ReadAll(resp.Body)
	assert.Nil(err)

	apiResp := new(ArchiveLineResponse)
	err = json.Unmarshal(body, apiResp)
	assert.Nil(err)
	assert.Equal("", apiResp.Error)

	b, err := ds.Get("history")
	assert.Nil(err)
	assert.True(regexp.MustCompile("^[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9] \t\tput water in pot$").Match(b))
}

func TestServer_ArchiveLine_LineMissing(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ts, cleanup := newBasicTaskstoreWithTestdata()
	defer cleanup()

	api := &Server{
		taskstore: ts,
	}
	addr := listenAddr()
	err := api.Start(addr)
	defer api.Stop()
	assert.Nil(err)

	resp, err := http.Get(
		"http://" + addr + "/archive_line/" + url.PathEscape(string(
			common.GetLineID("make_pasta", "flooptyboop"),
		)),
	)
	// Non-200 response does not cause an error on http.Get
	assert.Nil(err)
	assert.Equal(404, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.Nil(err)

	apiResp := new(ArchiveLineResponse)
	err = json.Unmarshal(body, apiResp)
	assert.Nil(err)
	assert.NotNil(apiResp.Error)
	assert.NotEqual("", apiResp.Error)
}

func TestServer_ArchiveLine_FileMissing(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ts, cleanup := newBasicTaskstoreWithTestdata()
	defer cleanup()

	api := &Server{
		taskstore: ts,
	}
	addr := listenAddr()
	err := api.Start(addr)
	defer api.Stop()
	assert.Nil(err)

	resp, err := http.Get(
		"http://" + addr + "/archive_line/" + url.PathEscape(string(
			common.GetLineID("nonexistent_file", "flooptyboop"),
		)),
	)
	// Non-200 response does not cause an error on http.Get
	assert.Nil(err)
	assert.Equal(404, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.Nil(err)

	apiResp := new(ArchiveLineResponse)
	err = json.Unmarshal(body, apiResp)
	assert.Nil(err)
	assert.NotNil(apiResp.Error)
	assert.NotEqual("", apiResp.Error)
}

func TestServer_InsertTask(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ds, cleanup := newFSDatastoreWithTestdata()
	defer cleanup()
	ts := NewBasicTaskstore(ds)

	api := &Server{
		taskstore: ts,
	}
	addr := listenAddr()
	err := api.Start(addr)
	defer api.Stop()
	assert.Nil(err)

	reqObj := &InsertTaskRequest{
		LineID: common.LineID("make_pasta:0"),
		Task:   common.NewTask(common.NewTreeNode("alpha")),
	}
	b, err := json.Marshal(reqObj)
	assert.Nil(err)

	resp, err := http.Post(
		"http://"+addr+"/insert_task/",
		"application/json",
		bytes.NewReader(b),
	)
	assert.Nil(err)
	body, err := io.ReadAll(resp.Body)
	assert.Nil(err)

	apiResp := new(InsertTaskResponse)
	err = json.Unmarshal(body, apiResp)
	assert.Nil(err)
	assert.Equal("", apiResp.Error)

	newFileBytes, err := ds.Get("make_pasta")
	assert.Nil(err)
	assert.True(regexp.MustCompile("alpha").Match(newFileBytes))
}
