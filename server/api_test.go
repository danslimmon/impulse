package server

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	//"github.com/danslimmon/impulse/common"
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

func TestGetTask(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ts, cleanup := NewBasicTaskstoreWithTestdata()
	defer cleanup()

	api := &Server{
		taskstore: ts,
	}
	addr := listenAddr()
	err := api.Start(addr)
	assert.Nil(err)

	apiReq := new(GetTaskRequest)
	apiResp := new(GetTaskResponse)
	err = api.GetTask(apiReq, apiResp)
	assert.Equal(nil, err)
	/*
		makePasta := common.MakePasta()[0]
			assert.Equal(1, len(apiResp.Tasks))
			assert.True(apiResp.Tasks[0].RootNode.Equal(makePasta.RootNode))
	*/
}
