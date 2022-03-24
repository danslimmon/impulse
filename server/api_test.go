package server

import (
	"testing"

	"github.com/danslimmon/impulse/common"
	"github.com/stretchr/testify/assert"
)

func TestGetTask(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	s, cleanup := NewServerWithTestdata()
	defer cleanup()

	// One task exactly
	apiReq := new(GetTaskRequest)
	apiReq.TaskIDs = []string{
		string(common.GetLineID("make_pasta", "make pasta")),
	}
	apiResp := new(GetTaskResponse)
	err := s.GetTask(apiReq, apiResp)
	assert.Equal(nil, err)

	makePasta := common.MakePasta()[0]
	assert.Equal(1, len(apiResp.Tasks))
	assert.True(apiResp.Tasks[0].RootNode.Equal(makePasta.RootNode))

	// Zero task IDs supplied
	apiReq = new(GetTaskRequest)
	apiReq.TaskIDs = []string{}
	apiResp = new(GetTaskResponse)
	err = s.GetTask(apiReq, apiResp)
	assert.Equal(nil, err)
	assert.Equal(0, len(apiResp.Tasks))

	// Multiple task IDs supplied
	apiReq = new(GetTaskRequest)
	apiReq.TaskIDs = []string{
		string(common.GetLineID("make_pasta", "make pasta")),
		string(common.GetLineID("multiple_nested", "task 0")),
		string(common.GetLineID("multiple_nested", "task 1")),
	}
	apiResp = new(GetTaskResponse)
	err = s.GetTask(apiReq, apiResp)
	assert.Equal(nil, err)

	multipleNested := common.MultipleNested()
	assert.Equal(3, len(apiResp.Tasks))
	assert.True(apiResp.Tasks[0].RootNode.Equal(makePasta.RootNode))
	assert.True(apiResp.Tasks[1].RootNode.Equal(multipleNested[0].RootNode))
	assert.True(apiResp.Tasks[2].RootNode.Equal(multipleNested[1].RootNode))

	// Same task ID supplied twice
	apiReq = new(GetTaskRequest)
	apiReq.TaskIDs = []string{
		string(common.GetLineID("multiple_nested", "task 1")),
		string(common.GetLineID("multiple_nested", "task 1")),
	}
	apiResp = new(GetTaskResponse)
	err = s.GetTask(apiReq, apiResp)
	assert.Equal(nil, err)

	assert.Equal(2, len(apiResp.Tasks))
	assert.True(apiResp.Tasks[0].RootNode.Equal(multipleNested[1].RootNode))
	assert.True(apiResp.Tasks[1].RootNode.Equal(multipleNested[1].RootNode))
}
