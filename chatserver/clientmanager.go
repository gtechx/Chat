package main

import (
	"fmt"
	"github.com/nature19862001/base/gtnet"
	"sync"
)

var clientaddrmap map[string]*Client
var clientuidmap map[uint64]*Client
var clientdelchan chan string
var clientaddchan chan *Client
var lastupdatecount int

var lock *sync.Mutex

func init() {
	clientaddrmap = make(map[string]*Client)
	clientuidmap = make(map[uint64]*Client)
	clientdelchan = make(chan string, 1024)
	clientaddchan = make(chan *Client, 1024)

	lock = new(sync.Mutex)
	go startClientOp()
}

func newClient(conn gtnet.IConn) *Client {
	c := &Client{conn: conn, lock: new(sync.Mutex), isVerifyed: false, clientAddr: conn.ConnAddr()}
	conn.SetMsgParser(c)
	conn.SetListener(c)
	go c.waitForLogin()
	clientaddchan <- c
	return c
}

func addUidMap(uid uint64, client *Client) {
	lock.Lock()
	defer lock.Unlock()

	clientuidmap[uid] = client
}

func removeUidMap(uid uint64) {
	lock.Lock()
	defer lock.Unlock()
	delete(clientuidmap, uid)
}

func removeClient(addr string) {
	clientdelchan <- addr
}

func sendMsgToUid(uid uint64, data []byte) {
	lock.Lock()
	defer lock.Unlock()

	client, ok := clientuidmap[uid]
	fmt.Println("send msg to uid:", uid)
	if ok {
		client.send(data)
	}
}

func startClientOp() {
	for {
		select {
		case newclient := <-clientaddchan:
			addr := newclient.conn.ConnAddr()
			clientaddrmap[addr] = newclient
		case deladdr := <-clientdelchan:
			delete(clientaddrmap, deladdr)
		}

		clientcount := len(clientaddrmap)
		deltacount := clientcount - lastupdatecount
		if deltacount >= 100 || deltacount <= -100 {
			gDataManager.incrServerClientCountBy(serverAddr, deltacount)
			lastupdatecount = clientcount
		}
	}
}
