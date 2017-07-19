package main

import (
	"flag"
	"fmt"
	. "github.com/nature19862001/Chat/common"
	"github.com/nature19862001/base/gtnet"
)

var server *gtnet.Server
var quit chan int

var nettype string = "tcp"
var addr string = "127.0.0.1:9090"
var redisnet string = "tcp"
var redisaddr string = "127.0.0.1:6379"

func main() {
	var err error

	pnet := flag.String("net", "tcp", "-net=")
	paddr := flag.String("addr", "127.0.0.1:9090", "-addr=")
	predisnet := flag.String("redisnet", "tcp", "-redisnet=")
	predisaddr := flag.String("redisaddr", "127.0.0.1:6379", "-redisaddr=")

	flag.Parse()

	nettype = *pnet
	addr = *paddr
	redisnet = *predisnet
	redisaddr = *predisaddr

	quit = make(chan int, 1)

	server = gtnet.NewServer(nettype, addr, onNewConn)

	//init redis
	err = InitRedis(redisnet, redisaddr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer CloseRedis()

	//start server
	err = server.Start()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer server.Stop()
	fmt.Println("server start ok...")

	<-quit
}

func onNewConn(conn gtnet.IConn) {
	addr := conn.ConnAddr()
	fmt.Println("new conn:", addr)
	newClient(conn)
}
