package server

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilesystemDatastore_Append(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ds, cleanup := newFSDatastoreWithTestdata()
	defer cleanup()

	// Append to a file that doesn't exist yet
	err := ds.Append("foo", []byte("first line\n"))
	assert.Nil(err)
	rslt, err := ds.Get("foo")
	assert.Nil(err)
	assert.Equal([]byte("first line\n"), rslt, fmt.Sprintf("unexpected file contents: '%s'", string(rslt)))

	// Append to a file that already exists
	err = ds.Append("foo", []byte("second line\n"))
	assert.Nil(err)
	rslt, err = ds.Get("foo")
	assert.Nil(err)
	assert.Equal([]byte("first line\nsecond line\n"), rslt, fmt.Sprintf("unexpected file contents: '%s'", string(rslt)))
}
