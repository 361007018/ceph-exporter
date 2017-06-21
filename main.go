package main

import (
	. "ceph-exporter/common"
	. "ceph-exporter/exporter"
	"errors"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"os"
)

var (
	threadCtl *ThreadCtl
	exporter  IExporter
)

func init() {
	// create thread controller
	threadCtl = CreateThreadCtl()

	// initial log
	initLog()

	// load configuration
	conf, err := config.NewConfig("ini", "conf/ceph-exporter.ini")
	if err != nil {
		logs.Error(err)
		threadCtl.MarkStopAsync()
		return
	}

	logs.Trace("src::host=" + conf.String("src::host"))
	logs.Trace("src::port=" + conf.String("src::port"))
	logs.Trace("dest::host=" + conf.String("dest::host"))
	logs.Trace("dest::port=" + conf.String("dest::port"))
	logs.Trace("dest::type=" + conf.String("dest::type"))

	// get source endpoint
	srcPort, err := conf.Int("src::port")
	if err != nil {
		logs.Error(err)
		threadCtl.MarkStopAsync()
		return
	}
	srcEndpoint := &Endpoint{
		Host: conf.String("src::host"),
		Port: int16(srcPort),
	}

	// get destination endpoint
	destPort, err := conf.Int("dest::port")
	if err != nil {
		logs.Error(err)
		threadCtl.MarkStopAsync()
		return
	}
	destEndpoint := &Endpoint{
		Host: conf.String("dest::host"),
		Port: int16(destPort),
	}

	// create exporter
	destType := conf.String("dest::type")
	switch destType {
	case "telegraf":
		{
			exporter = new(TelegrafExporter)
			exporter.Init(srcEndpoint, destEndpoint)
		}
	default:
		{
			logs.Error(errors.New("unsupport type of dest"))
			threadCtl.MarkStopAsync()
			return
		}
	}

	threadCtl.MarkStartAsync()
}

func main() {
	threadCtl.Run(func() {
		if exporter != nil {
			exporter.Run()
		} else {
			logs.Error("exporter undefined")
			threadCtl.MarkStopAsync()
		}
	})
}

// inital log
func initLog() {
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	logs.SetLogger(logs.AdapterConsole)
	os.MkdirAll("logs/", 0664)
	logs.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/ceph-exporter.log","maxlines":100000,"daily":true,"maxdays":7,"separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`)
	logs.Error("Oh,shit!!!It's error!!!")
}
