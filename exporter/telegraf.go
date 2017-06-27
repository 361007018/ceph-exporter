package exporter

import (
	"bytes"
	. "ceph-exporter/common"
	"errors"
	"github.com/astaxie/beego/logs"
	"net/http"
	"strconv"
	"time"
)

type TelegrafExporter struct {
	Exporter
	Db string
}

func (this *TelegrafExporter) Export(data string) {
	logs.Debug("export to telegraf:" + data)
	dataReader := bytes.NewReader([]byte(data))
	res, _ := http.Post(this.destEndpoint.ToString()+"/write?db="+this.Db, "plain/text", dataReader)
	logs.Debug(res.StatusCode)
}

func (this *TelegrafExporter) Init(srcEndpoint *Endpoint, destEndpoint *Endpoint, interval uint, args ...interface{}) {
	this.Exporter.Init(srcEndpoint, destEndpoint, interval)
	this.Db = args[0].(string)
}

func (this *TelegrafExporter) Run() {
	for {
		var data string
		// get cluster status
		logs.Info("Get cluster status.")
		clusterStatus, err := this.Exporter.GetClusterStatus()
		if err != nil {
			logs.Error(err)
		}
		if clusterStatus == nil {
			logs.Error(errors.New("Could not get cluster status."))
		} else {
			logs.Debug(clusterStatus)
			// get overall status
			data += "\n" + `ceph,type=cluster,property=overall_status msg="` + clusterStatus.Health.OverallStatus + `"`
			// get cluster storage space
			data += "\n" + `ceph,type=cluster,property=storage_space bytes_used=` + strconv.FormatUint(clusterStatus.Pgmap.BytesUsed, 10) + `,bytes_avail=` + strconv.FormatUint(clusterStatus.Pgmap.BytesAvail, 10) + `,bytes_total=` + strconv.FormatUint(clusterStatus.Pgmap.BytesTotal, 10) + `,used_percent=` + strconv.FormatFloat(float64(clusterStatus.Pgmap.BytesUsed)/float64(clusterStatus.Pgmap.BytesTotal), 'f', -1, 64) + `,avail_percent=` + strconv.FormatFloat(float64(clusterStatus.Pgmap.BytesAvail)/float64(clusterStatus.Pgmap.BytesTotal), 'f', -1, 64)
			// get detail
			for _, value := range clusterStatus.Health.Detail {
				data += "\n" + `ceph,type=cluster,property=summary msg="` + value + `"`
			}
			// get summary
			for _, value := range clusterStatus.Health.Summary {
				data += "\n" + `ceph,type=cluster,property=summary msg="` + value.Summary + `" ` + strconv.FormatInt(time.Now().UnixNano(), 10)
			}
			// get pgs by state
			for _, value := range clusterStatus.Pgmap.PgsByState {
				data += "\nceph,type=pg_state,state_name=" + value.StateName + " value=" + strconv.FormatUint(value.Count, 10)
			}
			// get client io rate
			data += "\nceph,type=pgmap,rw_rate_type=read_bytes_sec value=" + strconv.FormatUint(clusterStatus.Pgmap.ReadBytesSec, 10)
			data += "\nceph,type=pgmap,rw_rate_type=write_bytes_sec value=" + strconv.FormatUint(clusterStatus.Pgmap.WriteBytesSec, 10)
			// get client io ops
			data += "\nceph,type=pgmap,rw_ops_type=read_op_per_sec value=" + strconv.FormatUint(clusterStatus.Pgmap.ReadOpPerSec, 10)
			data += "\nceph,type=pgmap,rw_ops_type=write_op_per_sec value=" + strconv.FormatUint(clusterStatus.Pgmap.WriteOpPerSec, 10)
			// get mon summary
			data += "\nceph,type=mon,property=quorum_size up=" + strconv.Itoa(len(clusterStatus.Health.Timechecks.Mons)) + ",all=" + strconv.Itoa(len(clusterStatus.Monmap.Mons)) + `,percent=` + strconv.FormatFloat(float64(len(clusterStatus.Health.Timechecks.Mons))/float64(len(clusterStatus.Monmap.Mons)), 'f', -1, 64)
		}
		// get osd tree
		logs.Info("Get osd tree.")
		osdTree, err := this.Exporter.GetOsdTree()
		if err != nil {
			logs.Error(err)
		}
		if osdTree == nil {
			logs.Error(errors.New("Could not get osd tree"))
		} else {
			logs.Debug(osdTree)
			var osd_up_in, osd_up_out, osd_down_in, osd_down_out int64
			for _, value := range osdTree.Nodes {
				switch value.Type {
				case "osd":
					{
						data += "\nceph,name=" + value.Name + ",type=osd_tree,id=" + strconv.FormatInt(value.Id, 10) + `,status=` + value.Status + `,reweight=` + strconv.FormatFloat(value.Reweight, 'f', -1, 64) + " crush_weight=" + strconv.FormatFloat(value.CrushWeight, 'f', -1, 64)
						switch value.Status {
						case "up":
							{
								if value.Reweight != 0 {
									osd_up_in += 1
								} else {
									osd_up_out += 1
								}
							}
						case "down":
							{
								if value.Reweight != 0 {
									osd_down_in += 1
								} else {
									osd_down_out += 1
								}
							}
						}
					}
				}
			}
			data += "\nceph,type=osd,osd_state=up_and_in value=" + strconv.FormatInt(osd_up_in, 10)
			data += "\nceph,type=osd,osd_state=up_and_out value=" + strconv.FormatInt(osd_up_out, 10)
			data += "\nceph,type=osd,osd_state=down_and_in value=" + strconv.FormatInt(osd_down_in, 10)
			data += "\nceph,type=osd,osd_state=down_and_out value=" + strconv.FormatInt(osd_down_out, 10)
		}

		this.Export(data)
		time.Sleep(time.Duration(int(time.Second) * int(this.interval)))
	}
}
