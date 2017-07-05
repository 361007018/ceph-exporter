package collector

import (
	"libceph/common"
)

type Collector struct {
}

type ICollector interface {
	GetClusterStatus() (*common.ResStatus, error)
	GetOsdDf() (*common.ResOsdDf, error)
	GetOsdTree() (*common.ResOsdTree, error)
	GetPoolStats() (*common.ResPoolStats, error)
}
