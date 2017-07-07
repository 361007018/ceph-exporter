package main

import (
	"ceph-exporter/collecter"
	. "ceph-exporter/common"
	. "ceph-exporter/exporter"
	"errors"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"os"
	"time"
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
			for {
				logs.Info("Exporter start running")
				exporter.Run()
				logs.Info("Exporter stop running.Restart in 5 s.")
				time.Sleep(time.Duration(int(time.Second) * 5))
			}
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
	logs.Debug("src::protocol=" + conf.String("src::protocol"))
	logs.Debug("src::host=" + conf.String("src::host"))
	logs.Debug("src::port=" + conf.String("src::port"))
	logs.Debug("dest::protocol=" + conf.String("dest::protocol"))
	logs.Debug("dest::host=" + conf.String("dest::host"))
	logs.Debug("dest::port=" + conf.String("dest::port"))
	logs.Debug("dest::type=" + conf.String("dest::type"))

	// get source endpoint
	srcPort, err := conf.Int("src::port")
	if err != nil {
		return err
	}
	srcEndpoint := &Endpoint{
		Protocol: conf.String("src::protocol"),
		Host:     conf.String("src::host"),
		Port:     int16(srcPort),
	}

	// get destination endpoint
	destPort, err := conf.Int("dest::port")
	if err != nil {
		return err
	}
	destEndpoint := &Endpoint{
		Protocol: conf.String("dest::protocol"),
		Host:     conf.String("dest::host"),
		Port:     int16(destPort),
	}

	interval, err := conf.Int("service::interval")
	if err != nil {
		return err
	}

	collecter := new(collecter.CephRestAPICollecter)
	collecter.Init(srcEndpoint)
	// create exporter
	destType := conf.String("dest::type")
	switch destType {
	case "telegraf":
		{
			telegrafDb := conf.String("telegraf::db")
			logs.Debug("telegraf::db=" + telegrafDb)
			var telegrafExporter *TelegrafExporter = new(TelegrafExporter)
			telegrafExporter.Init(collecter, destEndpoint, uint(interval), telegrafDb)
			exporter = telegrafExporter
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
	logs.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/ceph-exporter.log","maxlines":100000,"daily":true,"maxdays":7,"separate":["error", "debug"]}`)
}
