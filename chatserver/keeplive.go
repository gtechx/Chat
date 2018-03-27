package main

import (
	//. "github.com/nature19862001/Chat/common"
	//"github.com/nature19862001/base/gtnet"
	"time"

	"github.com/nature19862001/Chat/chatserver/Data"
)

func keepLiveInit() {
	go startServerTTLKeep()
}

func startServerTTLKeep() {
	timer := time.NewTimer(time.Second * 30)

	select {
	case <-timer.C:
		cdata.Manager().SetServerTTL(serverAddr, 60)
	}
}
