package main

import (
	"github.com/nature19862001/base/gtnet"
	"sync"
)

var clientmap map[string]*Client
var clientdelchan chan string

func init() {
	clientmap = make(map[string]*Client, 0)
	clientdelchan = make(chan string, 1024)
}

func newClient(conn gtnet.IConn) *Client {
	c := &Client{conn: conn, lock: new(sync.Mutex), isVerifyed: false}
	conn.SetMsgParser(c)
	conn.SetListener(c)
	go c.waitForLogin()
	addNewClient(c)
	return c
}

func addNewClient(client *Client) {
	addr := client.conn.ConnAddr()
	clientmap[addr] = client
}

func removeClient(addr string) {
	clientdelchan <- addr
}

func startClientDel() {
	for addr := range clientdelchan {
		delete(clientmap, addr)
	}
}
