package main

import (
	"flag"
	"fmt"
	. "github.com/nature19862001/Chat/common"
	"github.com/nature19862001/base/gtnet"
)

func messagePullInit() {
	go startMessagePull()
}

func startMessagePull() {
	for {
		data := gDataManager.pullMsg(serverAddr, 15)

		if data != nil {
			//
		}
	}
}
