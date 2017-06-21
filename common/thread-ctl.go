package common

// Thread controller
type ThreadCtl struct {
	signalChan chan string
}

// Init function in thread controller
func (this *ThreadCtl) Init() {
	this.signalChan = make(chan string, 1)
}

// Mark a "start" signal to current thread controller
func (this *ThreadCtl) MarkStartAsync() {
	go func() {
		this.signalChan <- "start"
	}()
}

// Mark a "stop" signal to current thread controller
func (this *ThreadCtl) MarkStopAsync() {
	go func() {
		this.signalChan <- "stop"
	}()
}

// Running thread
func (this *ThreadCtl) Run(f func()) {
	for {
		signal := this.WaitSignal()
		switch signal {
		case "start":
			{
				go f()
			}
		case "stop":
			return
		}
	}
}

// Keep waiting until receive signal
func (this *ThreadCtl) WaitSignal() string {
	select {
	case signal := <-this.signalChan:
		return signal
	}
}

// ThreadCtl factory function
func CreateThreadCtl() *ThreadCtl {
	threadCtl := new(ThreadCtl)
	threadCtl.Init()
	return threadCtl
}
