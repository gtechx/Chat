package main

// import (
// 	"fmt"
// 	. "github.com/nature19862001/base/common"
// 	"github.com/nature19862001/base/gtnet"
// )

const MsgId_ReqLogin int16 = 1000

type ReqLogin struct {
	Uid      uint64
	Password [32]byte
}

const MsgId_ReqRetLogin int16 = 1001

type RetLogin struct {
	Result byte
}
