package api

import (
	"github.com/danslimmon/impulse/common"
)

type ListResponse struct {
	List *common.BlopList `json:"data"`
}
