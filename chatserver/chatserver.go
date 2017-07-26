package main

import (
	"fmt"
	. "github.com/nature19862001/Chat/common"
	"github.com/nature19862001/base/gtnet"
)

var server *gtnet.Server
var quit chan int

var nettype string = "tcp"
var addr string = "127.0.0.1:9090"

func chatServerInit(nettype, addr string) bool {
	var err error

	server = gtnet.NewServer(nettype, addr, onNewConn)

	err = server.Start()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer server.Stop()
	fmt.Println("server start ok...")

	return true
}

func onNewConn(conn gtnet.IConn) {
	addr := conn.ConnAddr()
	fmt.Println("new conn:", addr)
	newClient(conn)
}
