package main

import (
	. "ceph-exporter/common"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
)

var threadCtl *ThreadCtl

func init() {
	threadCtl = ThreadCtlInit()
	// load configuration
	conf, err := config.NewConfig("ini", "conf/ceph-exporter.ini")
	if err != nil {
		logs.Error(err)
		threadCtl.MarkStop()
	}

	logs.Trace("ceph-rest-api::host" + conf.String("ceph-rest-api::host"))
	logs.Trace("ceph-rest-api::port" + conf.String("ceph-rest-api::port"))
	logs.Trace("telegraf::host" + conf.String("telegraf::host"))
	logs.Trace("telegraf::port" + conf.String("telegraf::port"))
}

func main() {
	for {
		signal := threadCtl.WaitSignal()
		switch signal {
		case "stop":
			logs.Info("Receive stop signal.Exit.")
			return
		}
	}
}
