package main

import (
	"fmt"
	. "github.com/nature19862001/Chat/common"
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
}

func (this *Client) Close() {
	if this.conn != nil {
		if this.isVerifyed {
			n, err := RedisConn.Do("SREM", "user:online", String(this.uid))
			if err != nil {
				fmt.Println(err.Error())
			}
			if n == nil {
				fmt.Println("sadd cmd failed!")
			}
			n, err = RedisConn.Do("SREM", "user:online:password", this.password)
			if err != nil {
				fmt.Println(err.Error())
			}
			if n == nil {
				fmt.Println("sadd cmd failed!")
			}
		}
		removeClient(this.conn.ConnAddr())
		this.conn.Close()
		this.conn = nil
		this.isVerifyed = false
		this.state = state_none
	}
}

func (this *Client) waitForLogin() {
	this.state = state_connected
	this.timer = time.NewTimer(time.Second * 30)

	select {
	case <-this.timer.C:
		this.timer.Stop()
		this.lock.Lock()
		if !this.isVerifyed {
			this.Close()
			this.tickChan = make(chan int, 2)
			go this.startTick()
		}
		this.lock.Unlock()
	}
}

func (this *Client) startTick() {
	this.timer.Reset(time.Second * 30)
	for {
		select {
		case <-this.timer.C:
			this.countTimeOut++
			if this.countTimeOut >= 2 {
				this.timer.Stop()
				this.Close()
				return
			}
		case <-this.tickChan:
			this.timer.Reset(time.Second * 30)
		}
	}
}

func (this *Client) ParseHeader(data []byte) int {
	size := Int(data)
	//fmt.Println("header size :", size)
	//p.conn.Send(data)
	return size
}

func (this *Client) ParseMsg(data []byte) {
	//fmt.Println("client:", this.conn.ConnAddr(), "say:", String(data))
	msgid := Uint16(data)
	if this.isVerifyed {
		this.tickChan <- 1
	}
	switch msgid {
	case MsgId_ReqLogin:
		uid := Uint64(data[2:10])
		password := data[10:]
		ret := new(RetLogin)
		if checkLogin(uid, password) {
			this.state = state_logined
			this.uid = uid
			this.password = string(password)
			ret.Result = 1
			copy(ret.IP[0:], []byte("127.0.0.1"))
			ret.Port = 9090
			this.lock.Lock()
			this.isVerifyed = true
			this.lock.Unlock()
			n, err := RedisConn.Do("SADD", "user:online", String(uid))
			if err != nil {
				fmt.Println(err.Error())
			}
			if n == nil {
				fmt.Println("sadd cmd failed!")
			}
			n, err = RedisConn.Do("SADD", "user:online:password", string(password))
			if err != nil {
				fmt.Println(err.Error())
			}
			if n == nil {
				fmt.Println("sadd cmd failed!")
			}
		} else {
			ret.Result = 0
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
		ret := new(RetTick)
		ret.MsgId = MsgId_Tick
		this.send(Bytes(ret))
	case MsgId_ReqLoginOut:
		ret := new(RetLoginOut)
		ret.Result = 1
		ret.MsgId = MsgId_ReqRetLoginOut
		this.send(Bytes(ret))
		this.state = state_logouted
	default:
		fmt.Println("unknow msgid:", msgid)
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
	fmt.Println("tcpserver closed:", this.conn.ConnAddr())
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
