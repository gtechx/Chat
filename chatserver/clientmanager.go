package main

import (
	"github.com/nature19862001/base/gtnet"
	"sync"
)

var clientmap map[string]*Client
var clientdelchan chan string
var clientaddchan chan *Client
var lastupdatecount int

func init() {
	clientmap = make(map[string]*Client, 0)
	clientdelchan = make(chan string, 1024)
	clientaddchan = make(chan *Client, 1024)
	startClientOp()
}

func newClient(conn gtnet.IConn) *Client {
	c := &Client{conn: conn, lock: new(sync.Mutex), isVerifyed: false}
	conn.SetMsgParser(c)
	conn.SetListener(c)
	go c.waitForLogin()
	clientaddchan <- c
	return c
}

func removeClient(addr string) {
	clientdelchan <- addr
}

func startClientOp() {
	for {
		select {
		case newclient := clientaddchan:
			addr := newclient.conn.ConnAddr()
			clientmap[addr] = newclient
		case deladdr := clientdelchan:
			delete(clientmap, deladdr)
		}

		clientcount := len(clientmap)
		deltacount := clientcount - lastupdatecount
		if deltacount >= 100 || deltacount <= -100 {
			gDataManager.incrServerClientCountBy(serverAddr, deltacount)
			lastupdatecount = clientcount
		}
	}
}
