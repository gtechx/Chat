package main

import (
	"fmt"
	. "github.com/nature19862001/base/common"
	"github.com/nature19862001/base/gtnet"
)

type Processer struct {
	conn gtnet.IConn
}

func newProcesser(conn gtnet.IConn) *Processer {
	pro := &Processer{conn}
	conn.SetMsgParser(pro)
	conn.SetListener(pro)
	return pro
}

func (p *Processer) ParseHeader(data []byte) int {
	size := Int(data)
	//fmt.Println("header size :", size)
	//p.conn.Send(data)
	return size
}

func (p *Processer) ParseMsg(data []byte) {
	fmt.Println("client:", p.conn.ConnAddr(), "say:", String(data))
	p.conn.Send(append(Bytes(int16(len(data))), data...))
}

func (p *Processer) OnError(errorcode int, msg string) {
	fmt.Println("tcpserver error, errorcode:", errorcode, "msg:", msg)
}

func (p *Processer) OnPreSend([]byte) {

}

func (p *Processer) OnPostSend([]byte, int) {

}

func (p *Processer) OnClose() {
	fmt.Println("tcpserver closed:", p.conn.ConnAddr())
	//remove client conn
	delete(promap, p)
}

func (p *Processer) OnRecvBusy([]byte) {
	str := "server is busy"
	//p.conn.Send(Bytes(int16(len(str))))
	p.conn.Send(append(Bytes(int16(len(str))), []byte(str)...))
}

func (p *Processer) OnSendBusy([]byte) {
	// str := "server is busy"
	// p.conn.Send(Bytes(int16(len(str))))
	// p.conn.Send([]byte(str))
}
