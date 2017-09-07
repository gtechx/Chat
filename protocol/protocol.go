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
	MsgId_RetLogin

	MsgId_ReqAppLogin
	MsgId_RetAppLogin

	MsgId_ReqTokenLogin
	MsgId_RetTokenLogin

	MsgId_ReqToken
	MsgId_RetToken

	MsgId_ReqLoginOut
	MsgId_ReqRetLoginOut

	MsgId_ReqFriendList
	MsgId_RetFriendList

	MsgId_ReqFriendAdd
	MsgId_RetFriendAdd

	MsgId_FriendReqAgree
	MsgId_FriendReq
	MsgId_FriendReqResult

	MsgId_ReqFriendDel
	MsgId_RetFriendDel

	MsgId_ReqUserToBlack
	MsgId_RetUserToBlack

	MsgId_ReqRemoveUserInBlack
	MsgId_RetRemoveUserInBlack

	MsgId_ReqMoveFriendToGroup
	MsgId_RetMoveFriendToGroup

	MsgId_ReqSetFriendVerifyType
	MsgId_RetSetFriendVerifyType

	MsgId_Message
	MsgId_RetMessage

	MsgId_End
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

type MsgReqAppLogin struct {
	NullCmd
	Appname  [32]byte
	Password [32]byte
}

type MsgRetAppLogin struct {
	NullCmd
	Result uint16
}

type MsgReqTokenLogin struct {
	NullCmd
	Token []byte
}

type MsgRetTokenLogin struct {
	NullCmd
	Result uint16
	Count  byte
}

type MsgReqToken struct {
	NullCmd
	Uid      uint64
	Password [32]byte
}

type MsgRetToken struct {
	NullCmd
	Result uint16
	Uid    uint64
	Token  []byte
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

type MsgFriendReqResult struct {
	NullCmd
	Result uint16
	Fuid   uint64
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
	Fuid uint64
}

type MsgRetUserToBlack struct {
	NullCmd
	Result uint16
	Fuid   uint64
}

type MsgReqRemoveUserInBlack struct {
	NullCmd
	Fuid uint64
}

type MsgRetRemoveUserInBlack struct {
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

//app
const (
	AppMsgId_ReqTokenVerify uint16 = iota + MsgId_End
	AppMsgId_RetTokenVerify

	AppMsgId_ReqToken
	AppMsgId_RetToken
)

type AppMsgReqTokenVerify struct {
	NullCmd
	Token []byte
}

type AppMsgRetTokenVerify struct {
	NullCmd
	Result uint16
}

type AppMsgReqToken struct {
	NullCmd
	Uid uint64
}

type AppMsgRetToken struct {
	NullCmd
	Result uint16
	Token  []byte
}
