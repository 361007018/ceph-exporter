package exporter

import (
	. "ceph-exporter/common"
	"encoding/json"
	"io/ioutil"
	"libceph/common"
	"net/http"
)

type IExporter interface {
	Init(srcEndpoint *Endpoint, destEndpoint *Endpoint, interval uint, args ...interface{})
	Run()
}

// Base struct of exporter
type Exporter struct {
	srcEndpoint  *Endpoint
	destEndpoint *Endpoint
	interval     uint
}

// Initial exporter
func (this *Exporter) Init(srcEndpoint *Endpoint, destEndpoint *Endpoint, interval uint) {
	this.srcEndpoint = &Endpoint{
		Protocol: srcEndpoint.Protocol,
		Host:     srcEndpoint.Host,
		Port:     srcEndpoint.Port,
	}
	this.destEndpoint = &Endpoint{
		Protocol: destEndpoint.Protocol,
		Host:     destEndpoint.Host,
		Port:     destEndpoint.Port,
	}
	this.interval = interval
}

func (this *Exporter) GetClusterStatus() (*common.ResStatus, error) {
	httpClient := new(http.Client)
	resp, err := httpClient.Get(this.srcEndpoint.ToString() + "/v1/cluster/status")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := new(common.ResStatus)
	if err := json.Unmarshal(bytes, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *Exporter) GetOsdTree() (*common.ResOsdTree, error) {
	httpClient := new(http.Client)
	resp, err := httpClient.Get(this.srcEndpoint.ToString() + "/v1/osd/tree")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := new(common.ResOsdTree)
	if err := json.Unmarshal(bytes, result); err != nil {
		return nil, err
	}
	return result, nil
}
