package main

// import (
// 	"fmt"
// 	. "github.com/nature19862001/base/common"
// 	"github.com/nature19862001/base/gtnet"
// )

type NullCmd struct {
	MsgId uint16
}

const MsgId_ReqLogin uint16 = 1000

type ReqLogin struct {
	NullCmd
	Uid      uint64
	Password [32]byte
}

const MsgId_ReqRetLogin uint16 = 1001

type RetLogin struct {
	NullCmd
	Result byte
}

const MsgId_Tick uint16 = 1002

type RetTick struct {
	NullCmd
}

const MsgId_ReqLoginOut uint16 = 1003

type ReqLoginOut struct {
	NullCmd
}

const MsgId_ReqRetLoginOut uint16 = 1004

type RetLoginOut struct {
	NullCmd
	Result byte
}
