package main

import (
	"fmt"
	. "github.com/nature19862001/Chat/protocol"
	. "github.com/nature19862001/base/common"
	"github.com/nature19862001/base/gtnet"
	"time"
)

const (
	state_none          = 0
	state_connected int = 1
	state_logined   int = 2
	state_logouted  int = 3
	state_del       int = 4
	state_verify    int = 5
)

var msgProcesserMap map[uint16]func(*ChatClient, []byte)

func init() {
	msgProcesserMap = make(map[uint16]func(*ChatClient, []byte))
	msgProcesserMap[MsgId_ReqFriendList] = OnReqFriendList
	msgProcesserMap[MsgId_ReqFriendAdd] = OnReqFriendAdd
	msgProcesserMap[MsgId_ReqFriendDel] = OnReqFriendDel
	msgProcesserMap[MsgId_ReqUserToBlack] = OnReqUserToBlack
	msgProcesserMap[MsgId_ReqRemoveUserInBlack] = OnReqRemoveUserInBlack
	msgProcesserMap[MsgId_ReqMoveFriendToGroup] = OnReqMoveFriendToGroup
	msgProcesserMap[MsgId_ReqSetFriendVerifyType] = OnReqSetFriendVerifyType
	msgProcesserMap[MsgId_Message] = OnMessage
}

type ChatClient struct {
	conn gtnet.IConn
	uid  uint64
	//ChatClientAddr string

	recvChan chan []byte
}

func (this *ChatClient) Close() {
	fmt.Println("ChatClient:" + String(this.uid) + " closed")
	if this.conn != nil {
		this.conn.SetMsgParser(nil)
		this.conn.SetListener(nil)
		gDataManager.setUserOffline(this.uid)

		this.conn.Close()
		this.conn = nil
		close(this.recvChan)
		removeChatClient(this.uid)
	}
}

func (this *ChatClient) forceOffline() {
	//
}

func (this *ChatClient) serve() {
	this.recvChan = make(chan []byte, 2)
	this.conn.SetMsgParser(this)
	this.conn.SetListener(this)

	go this.startProcess()
}

func (this *ChatClient) startProcess() {
	timer := time.NewTimer(time.Second * 30)
	countTimeOut := 0

	for {
		select {
		case <-timer.C:
			fmt.Println("countTimeOut++")
			countTimeOut++
			if countTimeOut >= 2 {
				this.Close()
				goto end
			}
			timer.Reset(time.Second * 30)
		case data := <-this.recvChan:
			result := this.process(data)

			if result {
				goto end
			}

			countTimeOut = 0
			timer.Reset(time.Second * 30)
		}
	}
end:
	timer.Stop()
	this.Close()
	fmt.Println("tick end")
}

func (this *ChatClient) process(data []byte) bool {
	msgid := Uint16(data)

	switch msgid {
	case MsgId_Tick:
		ret := new(MsgTick)
		ret.MsgId = MsgId_Tick
		this.send(Bytes(ret))
	case MsgId_ReqLoginOut:
		ret := new(MsgRetLoginOut)
		ret.Result = 1
		ret.MsgId = MsgId_ReqRetLoginOut
		this.send(Bytes(ret))
		return true
	case MsgId_Echo:
		// ret := new(Echo)
		// ret.MsgId = MsgId_Echo
		// ret.Data = data[2:]
		this.send(data)
	default:
		fn, ok := msgProcesserMap[msgid]
		if ok {
			fn(this, data)
		} else {
			fmt.Println("unknown msgid:", msgid)
		}
	}

	return false
}

func (this *ChatClient) ParseHeader(data []byte) int {
	size := Int(data)
	//fmt.Println("header size :", size)
	//p.conn.Send(data)
	return size
}

func (this *ChatClient) ParseMsg(data []byte) {
	//fmt.Println("ChatClient:", this.conn.ConnAddr(), "say:", String(data))
	newdata := make([]byte, len(data))
	copy(newdata, data)
	this.recvChan <- newdata
}

func (this *ChatClient) send(buff []byte) {
	this.conn.Send(append(Bytes(int16(len(buff))), buff...))
}

func (this *ChatClient) OnError(errorcode int, msg string) {
	//fmt.Println("tcpserver error, errorcode:", errorcode, "msg:", msg)
}

func (this *ChatClient) OnPreSend([]byte) {

}

func (this *ChatClient) OnPostSend([]byte, int) {
	// if this.state == state_logouted {
	// 	this.Close()
	// }
}

func (this *ChatClient) OnClose() {
	//fmt.Println("tcpserver closed:", this.ChatClientAddr)
	this.Close()
}

func (this *ChatClient) OnRecvBusy([]byte) {
	//str := "server is busy"
	//p.conn.Send(Bytes(int16(len(str))))
	//this.conn.Send(append(Bytes(int16(len(str))), []byte(str)...))
}

func (this *ChatClient) OnSendBusy([]byte) {
	// str := "server is busy"
	// p.conn.Send(Bytes(int16(len(str))))
	// p.conn.Send([]byte(str))
}
