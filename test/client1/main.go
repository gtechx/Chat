package main

import (
	"flag"
	"fmt"
	. "github.com/nature19862001/Chat/protocol"
	. "github.com/nature19862001/base/common"
	"github.com/nature19862001/base/gtnet"
	"strings"
)

var client *gtnet.Client
var chatclient *gtnet.Client
var quit chan int

var nettype string = "tcp"
var addr string = "127.0.0.1:9090"
var user uint64 = 10001
var password string = "123456"

func main() {
	var err error

	pnet := flag.String("net", nettype, "-net=")
	paddr := flag.String("addr", addr, "-addr=")
	puser := flag.Uint64("u", user, "-u=")
	ppassword := flag.String("p", password, "-p=")

	flag.Parse()

	nettype = *pnet
	addr = *paddr
	user = *puser
	password = *ppassword

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

	req := new(MsgReqLogin)
	req.MsgId = MsgId_ReqLogin
	req.Uid = uint64(user)
	fmt.Println(Md5(password))
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

			//chatclient.Send(append(Bytes(int16(len(buff))), buff...))
			strarr := strings.Split(str, ":")
			cmd := strarr[0]
			data := ""
			if len(strarr) > 1 {
				data = strarr[1]
			}
			switch cmd {
			case "echo":
				req := new(MsgEcho)
				req.MsgId = MsgId_Echo
				req.Data = []byte(data)
				send(Bytes(req))
			case "flist":
				req := new(MsgReqFriendList)
				req.MsgId = MsgId_ReqFriendList
				send(Bytes(req))
			case "fadd":
				dataarr := strings.Split(data, "&")
				req := new(MsgReqFriendAdd)
				req.MsgId = MsgId_ReqFriendAdd
				req.Fuid = Uint64(dataarr[0])
				req.Group = []byte(dataarr[1])
				send(Bytes(req))
			case "msg":
				dataarr := strings.Split(data, "&")
				req := new(MsgMessage)
				req.MsgId = MsgId_Message
				req.Fuid = Uint64(dataarr[0])
				req.Data = []byte(dataarr[1])
				send(Bytes(req))
			case "black":
				req := new(MsgReqUserToBlack)
				req.MsgId = MsgId_ReqUserToBlack
				req.Fuid = Uint64(data)
				send(Bytes(req))
			case "rblack":
				req := new(MsgReqRemoveUserInBlack)
				req.MsgId = MsgId_ReqRemoveUserInBlack
				req.Fuid = Uint64(data)
				send(Bytes(req))
			case "movef":
				dataarr := strings.Split(data, "&")
				req := new(MsgReqMoveFriendToGroup)
				req.MsgId = MsgId_ReqMoveFriendToGroup
				req.Fuid = Uint64(dataarr[0])
				req.Group = []byte(dataarr[1])
				send(Bytes(req))
			case "setverify":
				req := new(MsgReqSetFriendVerifyType)
				req.MsgId = MsgId_ReqSetFriendVerifyType
				req.Type = Byte(data)
				send(Bytes(req))
			}
		}
	}
}
