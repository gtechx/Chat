package main

import (
	//. "github.com/nature19862001/Chat/common"
	//"github.com/nature19862001/base/gtnet"
	"time"
)

func serverMonitorInit() {
	go startServerMonitor()
}

func startServerMonitor() {
	timer := time.NewTimer(time.Second * 30)

	select {
	case <-timer.C:
		gDataManager.checkServerTTL()
	}
}
