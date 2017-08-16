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
	MsgId_RetFriendAdd

	MsgId_FriendReqAgree
	MsgId_FriendReq
	MsgId_FriendReqSuccess

	MsgId_ReqFriendDel
	MsgId_RetFriendDel

	MsgId_ReqUserToBlack
	MsgId_RetUserToBlack

	MsgId_ReqMoveFriendToGroup
	MsgId_RetMoveFriendToGroup

	MsgId_ReqSetFriendVerifyType
	MsgId_RetSetFriendVerifyType

	MsgId_Message
	MsgId_RetMessage
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
	Result uint16
	//IP     [16]byte
	//Port   uint16
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
	Result     uint16
	GroupCount byte
	Data       []byte
}

type MsgReqFriendAdd struct {
	NullCmd
	Fuid  uint64
	Group []byte
}

type MsgRetFriendAdd struct {
	NullCmd
	Result uint16
	Fuid   uint64
	Group  []byte
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

type MsgFriendReqSuccess struct {
	NullCmd
	Fuid uint64
}

type MsgReqFriendAddSuccess struct {
	NullCmd
}

type MsgReqFriendDel struct {
	NullCmd
	Fuid uint64
}

type MsgRetFriendDel struct {
	NullCmd
	Result uint16
	Fuid   uint64
}

type MsgReqUserToBlack struct {
	NullCmd
	Result uint16
	Fuid   uint64
}

type MsgRetUserToBlack struct {
	NullCmd
	Result uint16
	Fuid   uint64
}

type MsgReqMoveFriendToGroup struct {
	NullCmd
	Fuid  uint64
	Group []byte
}

type MsgRetMoveFriendToGroup struct {
	NullCmd
	Result uint16
	Fuid   uint64
}

type MsgReqSetFriendVerifyType struct {
	NullCmd
	Type byte
}

type MsgRetSetFriendVerifyType struct {
	NullCmd
	Result uint16
}

type MsgMessage struct {
	NullCmd
	Fuid uint64
	Data []byte
}

type MsgRetMessage struct {
	NullCmd
	Result uint16
	Fuid   uint64
}
