package main

import (
	"fmt"
	"github.com/nature19862001/base/gtnet"
	"sync"
)

var chatClientMap map[uint64]*ChatClient
var clientdelchan chan string
var clientaddchan chan *ChatClient

var chatlock *sync.Mutex
var lastupdatecount int

func init() {
	chatClientMap = make(map[uint64]*ChatClient)

	clientdelchan = make(chan string, 1024)
	clientaddchan = make(chan *ChatClient, 1024)

	chatlock = new(sync.Mutex)
}

func newChatClient(uid uint64, conn gtnet.IConn) *ChatClient {
	c := &ChatClient{conn: conn, uid: uid}

	chatlock.Lock()
	defer chatlock.Unlock()
	oldclient, ok := chatClientMap[uid]
	if ok {
		oldclient.forceOffline()
	}
	chatClientMap[uid] = c

	c.serve()
	fmt.Println("new chat client:", uid)

	clientcount := len(chatClientMap)
	deltacount := clientcount - lastupdatecount
	if deltacount >= 100 {
		gDataManager.incrServerClientCountBy(serverAddr, deltacount)
		lastupdatecount = clientcount
	}

	return c
}

func removeChatClient(uid uint64) {
	chatlock.Lock()
	defer chatlock.Unlock()

	delete(chatClientMap, uid)
	fmt.Println("delete chat client:", uid)

	clientcount := len(chatClientMap)
	deltacount := clientcount - lastupdatecount
	if deltacount <= -100 {
		gDataManager.incrServerClientCountBy(serverAddr, deltacount)
		lastupdatecount = clientcount
	}
}

func sendMsgToUid(uid uint64, data []byte) {
	chatlock.Lock()
	defer chatlock.Unlock()

	client, ok := chatClientMap[uid]
	fmt.Println("send msg to uid:", uid)
	if ok {
		client.send(data)
	} else {
		gDataManager.sendMsgToUser(uid, data)
	}
}

func cleanOnlineUsers() {
	chatlock.Lock()
	defer chatlock.Unlock()

	for uid, _ := range chatClientMap {
		gDataManager.setUserOffline(uid)
	}
	fmt.Println("cleanOnlineUsers end")
}

func verifyAppLogin(uid uint64) bool {
	// lock.Lock()
	// defer lock.Unlock()

	// client, ok := chatClientMap[uid]
	// if ok {
	// 	return client.verifyAppLogin()
	// } else {
	// 	return false
	// }
	return false
}

// func startChatClientOp() {
// 	for {
// 		select {
// 		case newclient := <-clientaddchan:
// 			addr := newclient.conn.ConnAddr()
// 			clientaddrmap[addr] = newclient
// 		case deladdr := <-clientdelchan:
// 			client, ok := clientaddrmap[deladdr]
// 			//if the same client close and login quickly and if with the same addr, we should check state to avoid delete a new client connect.
// 			if ok && client.state == state_del {
// 				delete(clientaddrmap, deladdr)
// 			}
// 		}

// 		clientcount := len(clientaddrmap)
// 		deltacount := clientcount - lastupdatecount
// 		if deltacount >= 100 || deltacount <= -100 {
// 			gDataManager.incrServerChatClientCountBy(serverAddr, deltacount)
// 			lastupdatecount = clientcount
// 		}
// 	}
// }
