package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/danslimmon/impulse/common"
	"github.com/danslimmon/impulse/server"
)

// Returns the path to a directory containing the contents of server/testdata/make_pasta as _.
//
// Also returns a cleanup function that should be called when the caller is done with the temporary
// directory.
func copyMakePastaDir() (string, func()) {
	d, err := ioutil.TempDir("", "impulse-test-*")
	if err != nil {
		panic(err.Error())
	}

	c := exec.Command("cp", "../server/testdata/make_pasta", path.Join(d, "_"))
	if err := c.Run(); err != nil {
		panic(err.Error())
	}

	cleanup := func() {
		os.RemoveAll(d)
	}

	return d, cleanup
}

// Returns a server/client pair with the server/testdata/make_pasta data as _.
//
// The returned server is already started.
//
// cleanup is a function that should be called when the test is done with the returned server and
// client.
func makePastaPair() (*server.ImpulseAPI, *ImpulseAPIClient, func()) {
	dataDir, copyCleanup := copyMakePastaDir()

	addr := "127.0.0.1:30272"
	ds := server.NewFilesystemDatastore(dataDir)
	ts := server.NewBasicTaskstore(ds)
	apiServer := server.NewImpulseAPI(ts)
	if apiServer.Start(addr) != nil {
		panic("failed to start test server on " + addr)
	}

	apiClient := NewImpulseAPIClient(addr)

	cleanup := func() {
		err := apiServer.Stop()
		if err != nil {
			panic(err.Error())
		}
		copyCleanup()
	}

	return apiServer, apiClient, cleanup
}

func Test_ImpulseAPIClient_GetTaskListResponse(t *testing.T) {
	// no t.Parallel() so we don't have to worry about giving out unique server ports
	assert := assert.New(t)

	_, client, cleanup := makePastaPair()
	defer cleanup()
	resp, err := client.GetTaskList("_")
	assert.Nil(err)
	assert.Equal(common.MakePasta(), resp.Result)
}
