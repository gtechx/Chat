package main

import (
	"fmt"
	. "github.com/nature19862001/Chat/protocol"
	. "github.com/nature19862001/base/common"
	"github.com/nature19862001/base/gtnet"
	"sync"
	"time"
)

const (
	state_none          = 0
	state_connected int = 1
	state_logined   int = 2
	state_logouted  int = 3
)

var msgProcesserMap map[uint16]func(*Client, []byte)

func init() {
	msgProcesserMap = make(map[uint16]func(*Client, []byte))
	msgProcesserMap[MsgId_ReqFriendList] = OnReqFriendList
	msgProcesserMap[MsgId_ReqFriendAdd] = OnReqFriendAdd
	msgProcesserMap[MsgId_ReqFriendDel] = OnReqFriendDel
	msgProcesserMap[MsgId_ReqUserToBlack] = OnReqUserToBlack
	msgProcesserMap[MsgId_ReqMoveFriendToGroup] = OnReqMoveFriendToGroup
	msgProcesserMap[MsgId_ReqSetFriendVerifyType] = OnReqSetFriendVerifyType
	msgProcesserMap[MsgId_Message] = OnMessage
}

type Client struct {
	conn         gtnet.IConn
	lock         *sync.Mutex
	isVerifyed   bool
	timer        *time.Timer
	verfiycount  int
	countTimeOut int
	tickChan     chan int
	uid          uint64
	password     string
	state        int
	clientAddr   string
}

func (this *Client) Close() {
	fmt.Println("client:" + this.clientAddr + "closed")
	if this.conn != nil {
		this.conn.SetMsgParser(nil)
		this.conn.SetListener(nil)
		if this.isVerifyed {
			gDataManager.setUserOffline(this.uid)
			// n, err := RedisConn.Do("SREM", "user:online", String(this.uid))
			// if err != nil {
			// 	fmt.Println(err.Error())
			// }
			// if n == nil {
			// 	fmt.Println("sadd cmd failed!")
			// }
			// n, err = RedisConn.Do("SREM", "user:online:password", this.password)
			// if err != nil {
			// 	fmt.Println(err.Error())
			// }
			// if n == nil {
			// 	fmt.Println("sadd cmd failed!")
			// }
		}
		removeClient(this.clientAddr)
		this.conn.Close()
		this.conn = nil
		this.isVerifyed = false
		this.state = state_none
	}

	this.closeTimer()
}

func (this *Client) closeTimer() {
	if this.timer != nil {
		this.timer.Reset(time.Millisecond * 1)
		this.timer = nil
	}
}

func (this *Client) waitForLogin() {
	this.state = state_connected
	this.timer = time.NewTimer(time.Second * 30)

	select {
	case <-this.timer.C:
		this.lock.Lock()
		if !this.isVerifyed {
			this.Close()
		}
		this.lock.Unlock()
	}
	fmt.Println("waitForLogin end")
}

func (this *Client) startTick() {
	this.timer = time.NewTimer(time.Second * 60)
	for {
		select {
		case <-this.timer.C:
			fmt.Println("countTimeOut++")
			this.countTimeOut++
			if this.countTimeOut >= 2 {
				if this.timer != nil {
					this.timer.Stop()
				}
				this.Close()
				return
			}
			if this.timer != nil {
				this.timer.Reset(time.Second * 60)
			}
		case <-this.tickChan:
			this.countTimeOut = 0
			if this.timer != nil {
				this.timer.Reset(time.Second * 60)
			}
		}

		if this.state == state_none {
			break
		}
	}
	fmt.Println("tick end")
}

func (this *Client) ParseHeader(data []byte) int {
	size := Int(data)
	fmt.Println("header size :", size)
	//p.conn.Send(data)
	return size
}

func (this *Client) ParseMsg(data []byte) {
	//fmt.Println("client:", this.conn.ConnAddr(), "say:", String(data))
	msgid := Uint16(data)
	fmt.Println("msgid:", msgid)
	if this.isVerifyed {
		this.tickChan <- 1
	} else if msgid != MsgId_ReqLogin {
		//if not logined, do not response to any msg
		return
	}

	switch msgid {
	case MsgId_ReqLogin:
		uid := Uint64(data[2:10])
		password := data[10:]
		ret := new(MsgRetLogin)
		code := gDataManager.checkLogin(uid, string(password))
		if code == ERR_NONE {
			this.state = state_logined
			this.uid = uid
			this.password = string(password)
			ret.Result = uint16(ERR_NONE)
			addUidMap(uid, this)
			//copy(ret.IP[0:], []byte("127.0.0.1"))
			//ret.Port = 9090
			ok := gDataManager.setUserOnline(uid)
			if !ok {
				ret.Result = uint16(ERR_REDIS)
			} else {
				this.lock.Lock()
				this.isVerifyed = true
				this.tickChan = make(chan int, 2)
				this.closeTimer()
				go this.startTick()
				this.lock.Unlock()
				fmt.Println("addr:" + this.conn.ConnAddr() + " logined success")
			}
		} else {
			ret.Result = uint16(code)
			this.verfiycount++

			if this.verfiycount < 5 {
				this.timer.Reset(time.Second * 30)
			} else {
				this.Close()
			}
		}
		ret.MsgId = MsgId_ReqRetLogin
		this.send(Bytes(ret))
	case MsgId_Tick:
		ret := new(MsgTick)
		ret.MsgId = MsgId_Tick
		this.send(Bytes(ret))
	case MsgId_ReqLoginOut:
		ret := new(MsgRetLoginOut)
		ret.Result = 1
		ret.MsgId = MsgId_ReqRetLoginOut
		this.send(Bytes(ret))
		this.state = state_logouted
	case MsgId_Echo:
		// ret := new(Echo)
		// ret.MsgId = MsgId_Echo
		// ret.Data = data[2:]
		this.send(data)
	default:
		fn, ok := msgProcesserMap[msgid]
		if ok {
			fn(this, data[2:])
		} else {
			fmt.Println("unknown msgid:", msgid)
		}
	}
}

func (this *Client) send(buff []byte) {
	this.conn.Send(append(Bytes(int16(len(buff))), buff...))
}

func (this *Client) OnError(errorcode int, msg string) {
	//fmt.Println("tcpserver error, errorcode:", errorcode, "msg:", msg)
}

func (this *Client) OnPreSend([]byte) {

}

func (this *Client) OnPostSend([]byte, int) {
	if this.state == state_logouted {
		this.Close()
	}
}

func (this *Client) OnClose() {
	//fmt.Println("tcpserver closed:", this.clientAddr)
	this.Close()
}

func (this *Client) OnRecvBusy([]byte) {
	//str := "server is busy"
	//p.conn.Send(Bytes(int16(len(str))))
	//this.conn.Send(append(Bytes(int16(len(str))), []byte(str)...))
}

func (this *Client) OnSendBusy([]byte) {
	// str := "server is busy"
	// p.conn.Send(Bytes(int16(len(str))))
	// p.conn.Send([]byte(str))
}
