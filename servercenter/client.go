package main

import (
	"fmt"
	. "github.com/nature19862001/Chat/common"
	. "github.com/nature19862001/base/common"
	"github.com/nature19862001/base/gtnet"
	"sync"
	"time"
)

type Client struct {
	conn        gtnet.IConn
	lock        *sync.Mutex
	isVerifyed  bool
	timer       *time.Timer
	verfiycount int
}

func (this *Client) Close() {
	if this.conn != nil {
		removeClient(this.conn.ConnAddr())
		this.conn.Close()
		this.conn = nil
	}
}

func (this *Client) waitForLogin() {
	this.timer = time.NewTimer(time.Second * 30)
	select {
	case <-this.timer.C:
		this.timer.Stop()
		this.lock.Lock()
		if !this.isVerifyed {
			this.Close()
		}
		this.lock.Unlock()
	}
}

func (this *Client) ParseHeader(data []byte) int {
	size := Int(data)
	//fmt.Println("header size :", size)
	//p.conn.Send(data)
	return size
}

func (this *Client) ParseMsg(data []byte) {
	fmt.Println("client:", this.conn.ConnAddr(), "say:", String(data))
	msgid := Int16(data)
	switch msgid {
	case MsgId_ReqLogin:
		uid := Uint64(data[2:10])
		password := data[10:]
		ret := new(RetLogin)
		if checkLogin(uid, password) {
			ret.Result = 1
			this.lock.Lock()
			this.isVerifyed = true
			this.lock.Unlock()
			n, err := RedisConn.Do("SADD", "user:online", String(uid))
			if err != nil {
				fmt.Println(err.Error())
			}
			if Int(n) <= 0 {
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
		retdata := Bytes(ret)
		this.conn.Send(append(Bytes(int16(len(retdata))), retdata...))
	default:
		fmt.Println("unknow msgid:", msgid)
	}
	//this.conn.Send(append(Bytes(int16(len(data))), data...))
}

func (this *Client) OnError(errorcode int, msg string) {
	fmt.Println("tcpserver error, errorcode:", errorcode, "msg:", msg)
}

func (this *Client) OnPreSend([]byte) {

}

func (this *Client) OnPostSend([]byte, int) {

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
