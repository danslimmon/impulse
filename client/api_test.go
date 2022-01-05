package client

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

	c := exec.Command("cp", "../server/testdata/make_pasta", path.Join(d, "make_pasta"))
	if err := c.Run(); err != nil {
		panic(err.Error())
	}

	cleanup := func() {
		os.RemoveAll(d)
	}

	return d, cleanup
}

// testServerAndClient returns a server/client pair with the server/testdata data.
//
// The returned server is already started.
//
// Also returns a function that should be called when the test is done with the returned server and
// client.
func testServerAndClient() (*server.Server, *Client, func()) {
	addr := "127.0.0.1:30272"
	ts, tsCleanup := server.NewBasicTaskstoreWithTestdata()

	apiServer := server.NewServer(ts)
	if err := apiServer.Start(addr); err != nil {
		panic("failed to start test server on " + addr + ": " + err.Error())
	}

	apiClient := NewClient(addr)

	cleanup := func() {
		apiServer.Stop()
		tsCleanup()
	}
	return apiServer, apiClient, cleanup
}

func Test_Client_GetTaskList(t *testing.T) {
	// no t.Parallel() so we don't have to worry about giving out unique server ports
	assert := assert.New(t)

	_, client, cleanup := testServerAndClient()
	defer cleanup()
	resp, err := client.GetTaskList("make_pasta")
	assert.Nil(err)
	assert.Equal(common.MakePasta(), resp.Result)
}

func Test_Client_GetTaskList_Nonexistent(t *testing.T) {
	// no t.Parallel() so we don't have to worry about giving out unique server ports
	assert := assert.New(t)

	_, client, cleanup := testServerAndClient()
	defer cleanup()
	_, err := client.GetTaskList("nonexistent_task_list")
	assert.NotNil(err)
}

func Test_Client_ArchiveLine(t *testing.T) {
	// no t.Parallel() so we don't have to worry about giving out unique server ports
	assert := assert.New(t)

	_, client, cleanup := testServerAndClient()
	defer cleanup()
	_, err := client.ArchiveLine(common.GetLineID("make_pasta", "\t\tput water in pot"))
	assert.Nil(err)
}

func Test_Client_ArchiveLine_Error(t *testing.T) {
	// no t.Parallel() so we don't have to worry about giving out unique server ports
	assert := assert.New(t)

	_, client, cleanup := testServerAndClient()
	defer cleanup()
	_, err := client.ArchiveLine("malformed_task_id")
	assert.NotNil(err)
}

func Test_Client_InsertTask(t *testing.T) {
	// no t.Parallel() so we don't have to worry about giving out unique server ports
	assert := assert.New(t)

	_, client, cleanup := testServerAndClient()
	defer cleanup()
	_, err := client.InsertTask(
		common.LineID("make_pasta:0"),
		common.NewTask(common.NewTreeNode("alpha")),
	)
	assert.Nil(err)
}
