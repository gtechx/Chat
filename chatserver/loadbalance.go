package main

import (
	"fmt"
	. "github.com/nature19862001/Chat/common"
	"github.com/nature19862001/base/gtnet"
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

}
