package exporter

import (
	"ceph-exporter/collecter"
	. "ceph-exporter/common"
)

type IExporter interface {
	Run()
}

// Base struct of exporter
type Exporter struct {
	srcEndpoint  *Endpoint
	destEndpoint *Endpoint
	interval     uint
	Collecter    collecter.ICollecter
}

// Initial exporter
func (this *Exporter) Init(collecter collecter.ICollecter, destEndpoint *Endpoint, interval uint) {
	this.Collecter = collecter
	this.destEndpoint = destEndpoint
	this.interval = interval
}
