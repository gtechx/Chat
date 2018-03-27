package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	//. "github.com/nature19862001/Chat/common"
	//"github.com/nature19862001/base/gtnet"
	"github.com/nature19862001/Chat/chatserver/Config"
	"github.com/nature19862001/Chat/chatserver/Data"
	"github.com/nature19862001/Chat/chatserver/Entity"
	"github.com/nature19862001/Chat/chatserver/Service"
	"github.com/nature19862001/base/gtnet"
)

var gDataManager dataManager
var quit chan os.Signal

var nettype string = "tcp"
var serverAddr string = "127.0.0.1:9090"
var redisNet string = "tcp"
var redisAddr string = "192.168.93.16:6379"

func main() {
	//var err error
	quit = make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	pnet := flag.String("net", nettype, "-net=")
	paddr := flag.String("addr", serverAddr, "-addr=")
	predisnet := flag.String("redisnet", redisNet, "-redisnet=")
	predisaddr := flag.String("redisaddr", redisAddr, "-redisaddr=")

	flag.Parse()

	if pnet != nil {
		config.ServerNet = *pnet
	}
	if paddr != nil {
		config.ServerAddr = *paddr
	}
	if predisnet != nil {
		config.RedisAddr = *predisnet
	}
	// nettype = *pnet
	// serverAddr = *paddr
	// redisNet = *predisnet
	// redisAddr = *predisaddr

	cdata.Manager().Initialize()
	centity.Manager().Initialize()

	//register server
	err := cdata.Manager().RegisterServer(config.ServerAddr)

	if err == nil {
		fmt.Println("register server to datamanager err:", err)
		return
	}

	//init loadbalance
	loadBanlanceInit()

	//init chat server
	// ok = chatServerStart(nettype, serverAddr)

	// if !ok {
	// 	fmt.Println("chat server init failed!!!")
	// 	return
	// }
	service := gtservice.NewService("chatserver", config.ServerNet, config.ServerAddr, onNewConn)
	err := service.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer service.Stop()

	//keep live init
	keepLiveInit()

	//other server live monitor init
	serverMonitorInit()

	//msg from other server monitor
	messagePullInit()

	<-quit

	//chatServerStop()
	centity.Manager().CleanOnlineUsers()
}

func onNewConn(conn gtnet.IConn) {
	centity.Manager().CreateNullEntity(conn)
}
