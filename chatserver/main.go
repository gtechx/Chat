package main

import (
	"flag"
	"fmt"
	//. "github.com/nature19862001/Chat/common"
	//"github.com/nature19862001/base/gtnet"
)

var gDataManager dataManager
var quit chan int

var nettype string = "tcp"
var serverAddr string = "127.0.0.1:9090"
var redisNet string = "tcp"
var redisAddr string = "127.0.0.1:6379"

func main() {
	//var err error

	pnet := flag.String("net", "tcp", "-net=")
	paddr := flag.String("addr", "127.0.0.1:9090", "-addr=")
	predisnet := flag.String("redisnet", "tcp", "-redisnet=")
	predisaddr := flag.String("redisaddr", "127.0.0.1:6379", "-redisaddr=")

	flag.Parse()

	nettype = *pnet
	serverAddr = *paddr
	redisNet = *predisnet
	redisAddr = *predisaddr

	quit = make(chan int, 1)
	gDataManager = new(redisDataManager)
	gDataManager.initialize()

	//register server
	ok := gDataManager.registerServer(serverAddr)

	if !ok {
		fmt.Println("can't register server to datamanager")
		return
	}

	//init loadbalance
	loadBanlanceInit()

	//init chat server
	ok = chatServerInit(nettype, serverAddr)

	if !ok {
		fmt.Println("chat server init failed!!!")
		return
	}

	//keep live init
	keepLiveInit()

	//other server live monitor init
	serverMonitorInit()

	//msg from other server monitor
	messagePullInit()

	<-quit
}
