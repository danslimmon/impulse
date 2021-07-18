package server

import (
	"testing"

	"github.com/danslimmon/impulse/common"
	"github.com/stretchr/testify/assert"
)

type MemFilesystem struct{}

func TestLoadBlopList_Empty(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	fs := &MemFilesystem{}

	exp := &common.BlopList{
		Name:  "_",
		Blops: []*Blop{},
	}

	rslt, err := loadBlopList(fs, "_")
	assert.Nil(err)
	assert.Equal(exp, rslt)
}
