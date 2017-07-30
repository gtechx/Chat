package main

import (
	"flag"
	"fmt"
	. "github.com/nature19862001/Chat/protocol"
	. "github.com/nature19862001/base/common"
	"github.com/nature19862001/base/gtnet"
)

var client *gtnet.Client
var chatclient *gtnet.Client
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
	client = gtnet.NewClient(nettype, addr)

	err = client.Connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer client.Close()
	newCenterClient(client)
	go startSend()

	req := new(ReqLogin)
	req.MsgId = MsgId_ReqLogin
	req.Uid = uint64(10001)
	fmt.Println(Md5("123456"))
	copy(req.Password[0:], []byte(Md5("123456")))
	send(Bytes(req))

	<-quit
}

func send(buff []byte) {
	client.Send(append(Bytes(int16(len(buff))), buff...))
}

func startSend() {
	var str string
	for {
		str = ""
		fmt.Scanln(&str)
		if str != "" {
			//bytes := Bytes(int16(len(str)))
			req := new(Echo)
			req.MsgId = MsgId_Echo
			req.Data = []byte(str)
			//fmt.Println(bytes)
			//buff := Bytes(req)
			send(Bytes(req))
			//chatclient.Send(append(Bytes(int16(len(buff))), buff...))
		}
	}
}
