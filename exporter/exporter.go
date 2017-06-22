package exporter

import (
	. "ceph-exporter/common"
	//. "libceph/common"
)

type IExporter interface {
	Init(srcEndpoint *Endpoint, destEndpoint *Endpoint, interval uint)
	Run()
}

// Base struct of exporter
type Exporter struct {
	srcEndpoint  *Endpoint
	destEndpoint *Endpoint
	threadCtl    *ThreadCtl
	interval     uint
}

// Initial exporter
func (this *Exporter) Init(srcEndpoint *Endpoint, destEndpoint *Endpoint, interval uint) {
	this.srcEndpoint = &Endpoint{
		Host: srcEndpoint.Host,
		Port: srcEndpoint.Port,
	}
	this.destEndpoint = &Endpoint{
		Host: destEndpoint.Host,
		Port: destEndpoint.Port,
	}
	this.interval = interval
	this.threadCtl = CreateThreadCtl()
}

func (this *Exporter) GetClusterStatus() {

}

func (this *Exporter) Run() {
}
