package exporter

import (
	. "ceph-exporter/common"
	//. "libceph/common"
)

type IExporter interface {
	Init(srcEndpoint *Endpoint, destEndpoint *Endpoint)
	Run()
}

// Base struct of exporter
type Exporter struct {
	srcEndpoint  *Endpoint
	destEndpoint *Endpoint
	threadCtl    *ThreadCtl
}

// Initial exporter
func (this *Exporter) Init(srcEndpoint *Endpoint, destEndpoint *Endpoint) {
	this.srcEndpoint = &Endpoint{
		Host: srcEndpoint.Host,
		Port: srcEndpoint.Port,
	}
	this.destEndpoint = &Endpoint{
		Host: destEndpoint.Host,
		Port: destEndpoint.Port,
	}
	this.threadCtl = CreateThreadCtl()
}

func (this *Exporter) GetClusterStatus() {

}

func (this *Exporter) Run() {
}
