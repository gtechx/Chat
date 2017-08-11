package protocol

// import (
// 	"fmt"
// 	. "github.com/nature19862001/base/common"
// 	"github.com/nature19862001/base/gtnet"
// )

const (
	MsgId_Tick uint16 = iota + 1000
	MsgId_Error
	MsgId_Echo

	MsgId_ReqLogin
	MsgId_ReqRetLogin

	MsgId_ReqLoginOut
	MsgId_ReqRetLoginOut
	MsgId_ReqFriendList
	MsgId_RetFriendList
	MsgId_ReqFriendAdd
	MsgId_FriendReqAgree
	MsgId_ReqFriendDel
	MsgId_RetFriendDel
	MsgId_ReqUserToBlack
	MsgId_RetUserToBlack
	MsgId_ReqMoveFriendToGroup
	MsgId_RetMoveFriendToGroup
	MsgId_ReqSetFriendVerifyType
	MsgId_RetSetFriendVerifyType
)

type NullCmd struct {
	MsgId uint16
}

type MsgError struct {
	NullCmd
	ErrorCode uint16
	ErrMsgId  uint16
}

func NewErrorMsg(errcode, errmsgid uint16) *MsgError {
	msg := new(MsgError)
	msg.MsgId = MsgId_Error
	msg.ErrorCode = errcode
	msg.ErrMsgId = errmsgid
	return new(MsgError)
}

type MsgTick struct {
	NullCmd
}

type MsgEcho struct {
	NullCmd
	Data []byte
}

type MsgReqLogin struct {
	NullCmd
	Uid      uint64
	Password [32]byte
}

type MsgRetLogin struct {
	NullCmd
	Result int8
	IP     [16]byte
	Port   uint16
}

type MsgReqLoginOut struct {
	NullCmd
}

type MsgRetLoginOut struct {
	NullCmd
	Result byte
}

type MsgReqFriendList struct {
	NullCmd
}

type MsgRetFriendList struct {
	NullCmd
	GroupCount byte
	Data       []byte
}

type MsgReqFriendAdd struct {
	NullCmd
	Fuid  uint64
	Group []byte
}

type MsgFriendReq struct {
	NullCmd
	Fuid uint64
}

type MsgFriendReqAgree struct {
	NullCmd
	Fuid  uint64
	Group []byte
}

type MsgReqFriendDel struct {
	NullCmd
	Fuid uint64
}

type MsgRetFriendDel struct {
	NullCmd
	Result byte
}

type MsgReqUserToBlack struct {
	NullCmd
	Fuid uint64
}

type MsgRetUserToBlack struct {
	NullCmd
	Result byte
}

type MsgReqMoveFriendToGroup struct {
	NullCmd
	Fuid  uint64
	Group []byte
}

type MsgRetMoveFriendToGroup struct {
	NullCmd
	Result byte
}

type MsgReqSetFriendVerifyType struct {
	NullCmd
	Type byte
}

type MsgRetSetFriendVerifyType struct {
	NullCmd
	Result byte
}
