package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/danslimmon/impulse/server"
)

// ImpulseAPIClient provides methods for using the Impulse API.
type ImpulseAPIClient struct {
	addr string
}

// url returns the full URL to the Impulse API endpoint with the given path.
func (apiClient *ImpulseAPIClient) url(path string) string {
	u := url.URL{
		Scheme: "http",
		Host:   apiClient.addr,
		Path:   path,
	}
	return u.String()
}

func (apiClient *ImpulseAPIClient) GetTaskList(listName string) (*server.GetTaskListResponse, error) {
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
		return respObj, fmt.Errorf("error: response code 404; body in logs")
	}
}

// NewImpulseAPIClient returns a fresh ImpulseAPIClient.
//
// addr is the host:port pair on which the server is listening.
func NewImpulseAPIClient(addr string) *ImpulseAPIClient {
	return &ImpulseAPIClient{addr: addr}
}
