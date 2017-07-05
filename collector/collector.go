package collector

import (
	"libceph/common"
)

type Collector struct {
}

type ICollector interface {
	GetClusterStatus() (*common.ResStatus, error)
	GetOsdTree() (*common.ResOsdTree, error)
	GetPoolStats() (*common.ResPoolStats, error)
}
