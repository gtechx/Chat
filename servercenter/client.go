package main

import (
	"fmt"
	. "github.com/nature19862001/base/common"
	"github.com/nature19862001/base/gtnet"
	"sync"
	"time"
)

type Client struct {
	conn       gtnet.IConn
	lock       *sync.Mutex
	isVerifyed bool
}

func newClient(conn gtnet.IConn) *Client {
	c := &Client{conn: conn, lock: new(sync.Mutex), isVerifyed: false}
	conn.SetMsgParser(c)
	conn.SetListener(c)
	go c.waitForLogin()
	return c
}

func (this *Client) Close() {
	if this.conn != nil {
		this.conn.Close()
		this.conn = nil
	}
}

func (this *Client) waitForLogin() {
	t1 := time.NewTimer(time.Second * 30)
	select {
	case <-t1.C:
		t1.Stop()
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
	msgid := int(Int16(data))
	switch msgid {
	case 1000:
	default:
		fmt.Println("unknow msgid:", msgid)
	}
	this.conn.Send(append(Bytes(int16(len(data))), data...))
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
	//remove client conn
	delete(promap, this)
}

func (this *Client) OnRecvBusy([]byte) {
	str := "server is busy"
	//p.conn.Send(Bytes(int16(len(str))))
	this.conn.Send(append(Bytes(int16(len(str))), []byte(str)...))
}

func (this *Client) OnSendBusy([]byte) {
	// str := "server is busy"
	// p.conn.Send(Bytes(int16(len(str))))
	// p.conn.Send([]byte(str))
}
