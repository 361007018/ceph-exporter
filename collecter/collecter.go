package collecter

import (
	"libceph/common"
)

type Collecter struct {
}

type ICollecter interface {
	GetClusterStatus() (*common.ResStatus, error)
	GetOsdDf() (*common.ResOsdDf, error)
	GetOsdTree() (*common.ResOsdTree, error)
	GetPoolStats() (*common.ResPoolStats, error)
}
