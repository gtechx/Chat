package main

import (
	"flag"
	"fmt"
	"github.com/nature19862001/base/gtnet"
)

var server *gtnet.Server
var promap map[*Processer]*Processer
var quit chan int

var nettype string = "tcp"
var addr string = "127.0.0.1:9090"

func main() {
	var err error

	pnet := flag.String("net", "tcp", "-net=")
	paddr := flag.String("addr", "127.0.0.1:9090", "-addr=")

	flag.Parse()

	nettype = *pnet
	addr = *paddr

	quit = make(chan int, 1)
	promap = make(map[*Processer]*Processer, 0)
	server = gtnet.NewServer(nettype, addr, onNewConn)

	err = server.Start()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("server start ok...")
	defer server.Stop()
	<-quit
}

func onNewConn(conn gtnet.IConn) {
	fmt.Println("new conn:", conn.ConnAddr())
	pro := newProcesser(conn)
	promap[pro] = pro
}
