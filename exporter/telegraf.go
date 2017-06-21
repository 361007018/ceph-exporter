package exporter

import (
	. "ceph-exporter/common"
)

type TelegrafExporter struct {
	Exporter
}

func (this *TelegrafExporter) Init(srcEndpoint *Endpoint, destEndpoint *Endpoint) {
	this.Exporter.Init(srcEndpoint, destEndpoint)
}

func (this *TelegrafExporter) Run() {

}
