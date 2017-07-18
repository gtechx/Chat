package main

import (
	"fmt"
	. "github.com/nature19862001/base/common"
	"github.com/nature19862001/base/gtnet"
)

type Client struct {
	conn gtnet.IConn
}

func newClient(conn gtnet.IConn) *Client {
	c := &Client{conn}
	conn.SetMsgParser(c)
	conn.SetListener(c)
	return c
}

func (this *Client) Close() {
	if this.conn != nil {
		this.conn.Close()
		this.conn = nil
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
