package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/danslimmon/impulse/common"
	"github.com/danslimmon/impulse/server"
)

// Client provides methods for using the Impulse API.
type Client struct {
	addr string
}

// url returns the full URL to the Impulse API endpoint with the given path.
func (apiClient *Client) url(path string) string {
	u := url.URL{
		Scheme: "http",
		Host:   apiClient.addr,
		Path:   path,
	}
	return u.String()
}

func (apiClient *Client) GetTaskList(listName string) (*server.GetTaskListResponse, error) {
	path := fmt.Sprintf("/tasklist/%s", listName)
	resp, err := http.Get(apiClient.url(path))
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	respObj := new(server.GetTaskListResponse)
	err = json.Unmarshal(b, respObj)
	if err != nil {
		return nil, err
	}
	if respObj.Error != "" {
		return nil, fmt.Errorf("Error response from server: %s", respObj.Error)
	}

	if resp.StatusCode == 200 {
		return respObj, nil
	} else {
		fmt.Printf("DEBUG: server response: '%s'", string(b))
		return respObj, fmt.Errorf("error: response code; body in logs")
	}
}

func (apiClient *Client) ArchiveLine(lineId common.LineID) (*server.ArchiveLineResponse, error) {
	path := fmt.Sprintf("/archive_line/%s", url.PathEscape(string(lineId)))
	resp, err := http.Get(apiClient.url(path))
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	respObj := new(server.ArchiveLineResponse)
	err = json.Unmarshal(b, respObj)
	if err != nil {
		return nil, err
	}
	if respObj.Error != "" {
		return nil, fmt.Errorf("Error response from server: %s", respObj.Error)
	}

	if resp.StatusCode == 200 {
		return respObj, nil
	} else {
		fmt.Printf("DEBUG: server response: '%s'", string(b))
		return respObj, fmt.Errorf("error: response code; body in logs")
	}
}

func (apiClient *Client) InsertTask(lineId common.LineID, task *common.Task) (*server.InsertTaskResponse, error) {
	reqObj := &server.InsertTaskRequest{
		LineID: lineId,
		Task:   task,
	}
	reqB, err := json.Marshal(reqObj)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal request object: %s", err.Error())
	}

	path := "/insert_task/"
	resp, err := http.Post(
		apiClient.url(path),
		"application/json",
		bytes.NewReader(reqB),
	)
	if err != nil {
		return nil, err
	}

	respB, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	respObj := new(server.InsertTaskResponse)
	err = json.Unmarshal(respB, respObj)
	if err != nil {
		return nil, err
	}
	if respObj.Error != "" {
		return nil, fmt.Errorf("Error response from server: %s", respObj.Error)
	}

	if resp.StatusCode == 200 {
		return respObj, nil
	} else {
		fmt.Printf("DEBUG: server response: '%s'", string(respB))
		return respObj, fmt.Errorf("error: response code; body in logs")
	}
}

// NewClient returns a fresh Client.
//
// addr is the host:port pair on which the server is listening.
func NewClient(addr string) *Client {
	return &Client{addr: addr}
}
