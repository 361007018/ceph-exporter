package exporter

import (
	. "ceph-exporter/common"
)

type TelegrafExporter struct {
	Exporter
}

func (this *TelegrafExporter) Init(srcEndpoint *Endpoint, destEndpoint *Endpoint, interval uint) {
	this.Exporter.Init(srcEndpoint, destEndpoint, interval)
}

func (this *TelegrafExporter) Run() {

}
