package main

import (
	"flag"
	"fmt"
	. "github.com/nature19862001/Chat/protocol"
	. "github.com/nature19862001/base/common"
	"github.com/nature19862001/base/gtnet"
)

var centerclient *gtnet.Client
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
	centerclient = gtnet.NewClient(nettype, addr)

	err = centerclient.Connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer centerclient.Close()
	newCenterClient(centerclient)
	go startSend()

	req := new(ReqLogin)
	req.MsgId = MsgId_ReqLogin
	req.Uid = uint64(1001)
	copy(req.Password[0:], []byte(Md5("1")))
	send(Bytes(req))

	<-quit
}

func send(buff []byte) {
	centerclient.Send(append(Bytes(int16(len(buff))), buff...))
}

func startSend() {
	var str string
	for {
		str = ""
		fmt.Scanln(&str)
		if str != "" {
			bytes := Bytes(int16(len(str)))
			//fmt.Println(bytes)
			chatclient.Send(append(bytes, []byte(str)...))
		}
	}
}
