package main

import (
	. "github.com/nature19862001/base/common"
)

var account map[uint64]string

func init() {
	account = make(map[uint64]string)
	for i := 1000; i < 10000; i++ {
		account[uint64(i)] = Md5("1")
	}
}

func checkLogin(uid uint64, password []byte) bool {
	pass, ok := account[uid]

	if ok && pass == string(password) {
		return true
	}

	return false
}
