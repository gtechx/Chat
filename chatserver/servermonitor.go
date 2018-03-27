package main

import (
	"time"

	"github.com/nature19862001/Chat/chatserver/Data"
)

func serverMonitorInit() {
	go startServerMonitor()
}

func startServerMonitor() {
	timer := time.NewTimer(time.Second * 30)

	select {
	case <-timer.C:
		cdata.Manager().CheckServerTTL()
	}
}
