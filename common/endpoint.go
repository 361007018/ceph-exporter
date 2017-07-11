package common

import (
	"strconv"
)

type Endpoint struct {
	Host     string
	Port     int16
	Protocol string
}

func (this *Endpoint) ToString() string {
	return this.Protocol + `://` + this.Host + `:` + strconv.Itoa(int(this.Port))
}
