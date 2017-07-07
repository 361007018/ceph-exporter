package collecter

import (
	. "ceph-exporter/common"
	"encoding/json"
	"io/ioutil"
	"libceph/common"
	"net/http"
)

type CephRestAPICollecter struct {
	srcEndpoint *Endpoint
}

func (this *CephRestAPICollecter) Init(endpoint *Endpoint) {
	this.srcEndpoint = endpoint
}

func (this *CephRestAPICollecter) GetClusterStatus() (*common.ResStatus, error) {
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

func (this *CephRestAPICollecter) GetOsdDf() (*common.ResOsdDf, error) {
	httpClient := new(http.Client)
	resp, err := httpClient.Get(this.srcEndpoint.ToString() + "/v1/osd/df")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := new(common.ResOsdDf)
	if err := json.Unmarshal(bytes, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (this *CephRestAPICollecter) GetOsdTree() (*common.ResOsdTree, error) {
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

func (this *CephRestAPICollecter) GetPoolStats() (*common.ResPoolStats, error) {
	httpClient := new(http.Client)
	resp, err := httpClient.Get(this.srcEndpoint.ToString() + "/v1/pool/stats/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := new(common.ResPoolStats)
	if err := json.Unmarshal(bytes, result); err != nil {
		return nil, err
	}
	return result, nil
}
