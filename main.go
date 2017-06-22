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

	// initial exporter
	if err := initExporter(); err != nil {
		logs.Error(err)
		threadCtl.MarkStopAsync()
		return
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

// initial exporter
func initExporter() error {
	// load configuration
	conf, err := config.NewConfig("ini", "conf/ceph-exporter.ini")
	if err != nil {
		return err
	}
	logs.Debug("service::interval=" + conf.String("service::interval"))
	logs.Debug("src::host=" + conf.String("src::host"))
	logs.Debug("src::port=" + conf.String("src::port"))
	logs.Debug("dest::host=" + conf.String("dest::host"))
	logs.Debug("dest::port=" + conf.String("dest::port"))
	logs.Debug("dest::type=" + conf.String("dest::type"))

	// get source endpoint
	srcPort, err := conf.Int("src::port")
	if err != nil {
		return err
	}
	srcEndpoint := &Endpoint{
		Host: conf.String("src::host"),
		Port: int16(srcPort),
	}

	// get destination endpoint
	destPort, err := conf.Int("dest::port")
	if err != nil {
		return err
	}
	destEndpoint := &Endpoint{
		Host: conf.String("dest::host"),
		Port: int16(destPort),
	}

	interval, err := conf.Int("service::interval")
	if err != nil {
		return err
	}

	// create exporter
	destType := conf.String("dest::type")
	switch destType {
	case "telegraf":
		{
			exporter = new(TelegrafExporter)
			exporter.Init(srcEndpoint, destEndpoint, uint(interval))
		}
	default:
		{
			return errors.New("unsupport type of dest")
		}
	}
	return nil
}

// initial log
func initLog() {
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	logs.SetLogger(logs.AdapterConsole)
	os.MkdirAll("logs/", 0664)
	logs.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/ceph-exporter.log","maxlines":100000,"daily":true,"maxdays":7,"separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`)
}
