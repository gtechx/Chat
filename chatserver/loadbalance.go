package main

import (
	//"fmt"
	//. "github.com/nature19862001/Chat/common"
	//"github.com/nature19862001/base/gtnet"
	"io"
	"net/http"
)

func loadBanlanceInit() {
	go startHTTPServer()
}

func startHTTPServer() {
	http.HandleFunc("/serverlist", listCmd)
	http.ListenAndServe(":9001", nil)
}

func getServerList(rw http.ResponseWriter, req *http.Request) {
	serverlist := gDataManager.getServerList()

	ret := "{\r\n\tserverlist:\r\n\t[\r\n"
	for i := 0; i < len(serverlist); i++ {
		ret += "\t\t{ addr:\"" + serverlist[i] + "\" },\r\n"
	}
	ret += "\t]\r\n"
	ret += "}\r\n"

	io.WriteString(rw, ret)
}
