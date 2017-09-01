package main

import (
	"github.com/nature19862001/base/gtnet"
	"sync"
)

var appClientMap map[string][]*AppClient
var lock *sync.Mutex

func init() {
	appClientMap = make(map[string][]*AppClient)
	lock = new(sync.Mutex)
}

func newAppClient(appname string, conn gtnet.IConn) *AppClient {
	c := &AppClient{conn: conn, appName: appname}

	lock.Lock()
	defer lock.Unlock()
	arr, ok := appClientMap[appname]
	if ok {
		appClientMap[appname] = append(arr, c)
	} else {
		appClientMap[appname] = []*AppClient{c}
	}

	c.serve()

	return c
}

func removeAppClient(appname string) {
	lock.Lock()
	defer lock.Unlock()
	delete(appClientMap, appname)
}
