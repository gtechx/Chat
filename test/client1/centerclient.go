package main

import (
	"fmt"
	. "github.com/nature19862001/Chat/protocol"
	. "github.com/nature19862001/base/common"
	"github.com/nature19862001/base/gtnet"
	"time"
)

type CenterClient struct {
	client    *gtnet.Client
	isLogined bool
}

func newCenterClient(client *gtnet.Client) *CenterClient {
	pro := &CenterClient{client: client}
	client.SetMsgParser(pro)
	client.SetListener(pro)
	return pro
}

func (this *CenterClient) startTick() {
	timer := time.NewTimer(time.Second * 30)
	for {
		select {
		case <-timer.C:
			//fmt.Println("send tick to server")
			req := new(MsgTick)
			req.MsgId = MsgId_Tick
			this.send(Bytes(req))
			timer.Reset(time.Second * 30)
		}
	}
}

func (this *CenterClient) ParseHeader(data []byte) int {
	size := Int(data)
	//fmt.Println("header size :", size)
	return size
}

func (this *CenterClient) send(buff []byte) {
	this.client.Send(append(Bytes(int16(len(buff))), buff...))
}

func (this *CenterClient) ParseMsg(data []byte) {
	msgid := Uint16(data)
	//fmt.Println("msgid:", msgid)
	switch msgid {
	case MsgId_ReqRetLogin:
		result := Uint16(data[2:4])

		if result == 0 {
			this.isLogined = true
			fmt.Println("login to server center success!")
			go this.startTick()
		} else {
			fmt.Println("login to server center failed! errcode:", result)
		}
	case MsgId_Tick:
		//fmt.Println("recv tick rom server")
	case MsgId_ReqLoginOut:
	case MsgId_Echo:
		fmt.Println("recv:" + String(data[2:]))
	case MsgId_RetFriendList:
		//
		fmt.Println("MsgId_RetFriendList result:", Uint16(data[2:4]), " groupcount:", data[4:5])
	case MsgId_RetFriendAdd:
		fmt.Println("friend add result:", Uint16(data[2:4]), Uint64(data[4:12]))
	case MsgId_FriendReqResult:
		fmt.Println("friend add req sended result:", Uint16(data[2:4]), Uint64(data[4:12]))
	case MsgId_FriendReqAgree:
		fmt.Println("friend had agreed your friend req:", Uint64(data[2:10]), string(data[10:]))
	case MsgId_FriendReq:
		fmt.Println("recv friend req:", Uint64(data[2:]))
		req := new(MsgReqFriendAdd)
		req.MsgId = MsgId_ReqFriendAdd
		req.Fuid = Uint64(data[2:])
		req.Group = []byte("aaa")
		this.send(Bytes(req))
	case MsgId_Message:
		fmt.Println("msg from", Uint64(data[2:10]), string(data[10:]))
	case MsgId_RetMessage:
		fmt.Println("msg send result:", Uint16(data[2:4]), Uint64(data[4:]))
	case MsgId_RetUserToBlack:
		fmt.Println("result ", Uint16(data[2:4]), " for user to black uid:", Uint64(data[4:]))
	case MsgId_RetRemoveUserInBlack:
		fmt.Println("result ", Uint16(data[2:4]), " for remove user in black uid:", Uint64(data[4:]))
	case MsgId_RetFriendDel:
		fmt.Println("result ", Uint16(data[2:4]), " for friend del uid:", Uint64(data[4:]))
	case MsgId_RetMoveFriendToGroup:
		fmt.Println("result ", Uint16(data[2:4]), "for move friend to group, uid:", Uint64(data[4:]))
	case MsgId_RetSetFriendVerifyType:
		fmt.Println("result ", Uint16(data[2:4]), "for set friend verify type")
	default:
		fmt.Println("unknow msgid:", msgid)
	}
}

func (this *CenterClient) OnError(errorcode int, msg string) {
	fmt.Println("tcpclient error, errorcode:", errorcode, "msg:", msg)
}

func (this *CenterClient) OnPreSend([]byte) {

}

func (this *CenterClient) OnPostSend([]byte, int) {

}

func (this *CenterClient) OnClose() {
	fmt.Println("tcpclient closed")
	quit <- 1
}

func (this *CenterClient) OnRecvBusy(buff []byte) {
	fmt.Println("client is busy for recv, msg size is ", len(buff))
}

func (this *CenterClient) OnSendBusy(buff []byte) {
	fmt.Println("client is busy for send, msg size is ", len(buff))
}
