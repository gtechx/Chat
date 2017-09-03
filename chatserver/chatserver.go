package main

import (
	"fmt"
	"github.com/nature19862001/base/gtnet"
)

var server *gtnet.Server

func chatServerStart(nettype, addr string) bool {
	var err error

	server = gtnet.NewServer(nettype, addr, onNewConn)

	err = server.Start()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	fmt.Println("server start ok...")

	return true
}

func chatServerStop() {
	if server != nil {
		server.Stop()
		fmt.Println("server stoped...")
	}
}

func onNewConn(conn gtnet.IConn) {
	addr := conn.ConnAddr()
	fmt.Println("new conn:", addr)
	newClient(conn)
}
