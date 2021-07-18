package server

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/danslimmon/impulse/common"
)

type Filesystem interface {
	GetBlopList(string) (*common.BlopList, error)
	GetFileContents(string) ([]byte, error)
}

type DiskFilesystem struct {
	rootDir string
}

func (fs *DiskFilesystem) absPath(filename string) string {
	return filepath.Join(fs.rootDir, filename)
}

func (fs *DiskFilesystem) GetFileContents(filename string) ([]byte, error) {
	return ioutil.ReadFile(fs.absPath(filename))
}
