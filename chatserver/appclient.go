package main

import (
	"fmt"
	. "github.com/nature19862001/Chat/protocol"
	. "github.com/nature19862001/base/common"
	"github.com/nature19862001/base/gtnet"
	"time"
)

var appMsgProcesserMap map[uint16]func(*AppClient, []byte)

func init() {
	appMsgProcesserMap = make(map[uint16]func(*AppClient, []byte))
	// appMsgProcesserMap[MsgId_ReqFriendList] = OnReqFriendList
	// appMsgProcesserMap[MsgId_ReqFriendAdd] = OnReqFriendAdd
	// appMsgProcesserMap[MsgId_ReqFriendDel] = OnReqFriendDel
	// appMsgProcesserMap[MsgId_ReqUserToBlack] = OnReqUserToBlack
	// appMsgProcesserMap[MsgId_ReqRemoveUserInBlack] = OnReqRemoveUserInBlack
	// appMsgProcesserMap[MsgId_ReqMoveFriendToGroup] = OnReqMoveFriendToGroup
	// appMsgProcesserMap[MsgId_ReqSetFriendVerifyType] = OnReqSetFriendVerifyType
	// appMsgProcesserMap[MsgId_Message] = OnMessage
}

type AppClient struct {
	conn     gtnet.IConn
	appName  string
	recvChan chan []byte
}

func (this *AppClient) Close() {
	fmt.Println("AppClient:" + this.appName + " closed")
	if this.conn != nil {
		this.conn.SetMsgParser(nil)
		this.conn.SetListener(nil)
		gDataManager.setAppOffline(this.appName)

		this.conn.Close()
		this.conn = nil
		close(this.recvChan)
		removeAppClient(this.appName)
	}
}

func (this *AppClient) serve() {
	this.recvChan = make(chan []byte, 2)
	this.conn.SetMsgParser(this)
	this.conn.SetListener(this)

	go this.startProcess()
}

func (this *AppClient) startProcess() {
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

func (this *AppClient) process(data []byte) bool {
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
		fn, ok := appMsgProcesserMap[msgid]
		if ok {
			fn(this, data)
		} else {
			fmt.Println("unknown msgid:", msgid)
		}
	}

	return false
}

func (this *AppClient) ParseHeader(data []byte) int {
	size := Int(data)
	//fmt.Println("header size :", size)
	return size
}

func (this *AppClient) ParseMsg(data []byte) {
	//fmt.Println("AppClient:", this.conn.ConnAddr(), "say:", String(data))
	newdata := make([]byte, len(data))
	copy(newdata, data)
	this.recvChan <- newdata
}

func (this *AppClient) send(buff []byte) {
	this.conn.Send(append(Bytes(int16(len(buff))), buff...))
}

func (this *AppClient) OnError(errorcode int, msg string) {
	//fmt.Println("tcpserver error, errorcode:", errorcode, "msg:", msg)
}

func (this *AppClient) OnPreSend([]byte) {

}

func (this *AppClient) OnPostSend([]byte, int) {
	// if this.state == state_logouted {
	// 	this.Close()
	// }
}

func (this *AppClient) OnClose() {
	//fmt.Println("tcpserver closed:", this.AppClientAddr)
	this.Close()
}

func (this *AppClient) OnRecvBusy([]byte) {
	//str := "server is busy"
	//p.conn.Send(Bytes(int16(len(str))))
	//this.conn.Send(append(Bytes(int16(len(str))), []byte(str)...))
}

func (this *AppClient) OnSendBusy([]byte) {
	// str := "server is busy"
	// p.conn.Send(Bytes(int16(len(str))))
	// p.conn.Send([]byte(str))
}
