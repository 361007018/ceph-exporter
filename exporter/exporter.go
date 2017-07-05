package exporter

import (
	"ceph_exporter/collector"
	. "ceph_exporter/common"
)

type IExporter interface {
	Run()
}

// Base struct of exporter
type Exporter struct {
	srcEndpoint  *Endpoint
	destEndpoint *Endpoint
	interval     uint
	Collector    collector.ICollector
}

// Initial exporter
func (this *Exporter) Init(collector collector.ICollector, destEndpoint *Endpoint, interval uint) {
	this.Collector = collector
	this.destEndpoint = destEndpoint
	this.interval = interval
}
