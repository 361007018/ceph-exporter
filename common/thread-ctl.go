package common

// Thread controller
type ThreadCtl struct {
	signalChan chan string
}

// Init function in thread controller
func (this *ThreadCtl) Init() {
	this.signalChan = make(chan string, 1)
}

// Mark a "stop" signal to current thread controller
func (this *ThreadCtl) MarkStop() {
	go func() {
		this.signalChan <- "stop"
	}()
}

// Keep waiting until receive signal
func (this *ThreadCtl) WaitSignal() string {
	select {
	case signal := <-this.signalChan:
		return signal
	}
}

// ThreadCtl factory function
func ThreadCtlInit() *ThreadCtl {
	threadCtl := new(ThreadCtl)
	threadCtl.Init()
	return threadCtl
}
