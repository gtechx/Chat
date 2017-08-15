package main

import (
	. "github.com/nature19862001/Chat/protocol"
	. "github.com/nature19862001/base/common"
)

func OnReqFriendList(client *Client, data []byte) {
	flist, ret := gDataManager.getFriendList(client.uid)

	// if err != nil {
	// 	reterr := NewErrorMsg(ERR_REDIS, MsgId_ReqFriendList)
	// 	client.send(Bytes(reterr))
	// 	return
	// }

	retmsg := new(MsgRetFriendList)
	retmsg.MsgId = MsgId_RetFriendList
	retmsg.Result = uint16(ret)
	retmsg.GroupCount = byte(len(flist))

	data = []byte{}

	for groupname, uidarr := range flist {
		data = append(data, byte(len(groupname)))
		data = append(data, []byte(groupname)...)
		data = append(data, Bytes(uint16(len(uidarr)))...)

		for _, uid := range uidarr {
			data = append(data, Bytes(uid)...)
		}
	}

	retmsg.Data = data

	client.send(Bytes(retmsg))
}

func OnReqFriendAdd(client *Client, data []byte) {
	fuid := Uint64(data)
	group := data[8:]
	ret := gDataManager.addFriend(client.uid, fuid, string(group))

	switch ret {
	case ERR_FRIEND_IN_BLACKLIST:
		//do nothing
	case ERR_REDIS:
		reterr := NewErrorMsg(ERR_REDIS, MsgId_ReqFriendAdd)
		client.send(Bytes(reterr))
	case ERR_FRIEND_EXIST:
		reterr := NewErrorMsg(ERR_FRIEND_EXIST, MsgId_ReqFriendAdd)
		client.send(Bytes(reterr))
	case ERR_USER_NOT_EXIST:
		reterr := NewErrorMsg(ERR_USER_NOT_EXIST, MsgId_ReqFriendAdd)
		client.send(Bytes(reterr))
	case ERR_FRIEND_MAX:
		reterr := NewErrorMsg(ERR_FRIEND_MAX, MsgId_ReqFriendAdd)
		client.send(Bytes(reterr))
	case ERR_FRIEND_ADD_NEED_REQ:
		req := new(MsgFriendReq)
		req.MsgId = MsgId_FriendReq
		req.Fuid = client.uid
		code := gDataManager.sendMsgToUser(fuid, Bytes(req))

		if code != ERR_NONE {
			reterr := NewErrorMsg(ERR_REDIS, MsgId_ReqFriendAdd)
			client.send(Bytes(reterr))
			return
		}

		//send req success msg to client
		retmsg := new(MsgFriendReqSuccess)
		retmsg.MsgId = MsgId_FriendReqSuccess
		client.send(Bytes(retmsg))
	case ERR_NONE:
		//send add success msg to client
		retmsg := new(MsgReqFriendAddSuccess)
		retmsg.MsgId = MsgId_ReqFriendAddSuccess
		client.send(Bytes(retmsg))
	}
}

func OnReqFriendDel(client *Client, data []byte) {
	fuid := Uint64(data)
	ret := gDataManager.deleteFriend(client.uid, fuid)

	switch ret {
	case ERR_FRIEND_NOT_EXIST:
		reterr := NewErrorMsg(ERR_FRIEND_NOT_EXIST, MsgId_ReqFriendDel)
		client.send(Bytes(reterr))
	case ERR_REDIS:
		reterr := NewErrorMsg(ERR_REDIS, MsgId_ReqFriendDel)
		client.send(Bytes(reterr))
	case ERR_NONE:
		//send add success msg to client
		retmsg := new(MsgFriendDelSucess)
		retmsg.MsgId = MsgId_FriendDelSuccess
		client.send(Bytes(retmsg))
	}
}

func OnReqUserToBlack(client *Client, data []byte) {
	fuid := Uint64(data)
	ret := gDataManager.addUserToBlacklist(client.uid, fuid)

	switch ret {
	case ERR_REDIS:
		retmsg := new(MsgRetUserToBlack)
		retmsg.MsgId = MsgId_RetUserToBlack
		retmsg.Result = uint16(ERR_REDIS)
		retmsg.Fuid = fuid
		client.send(Bytes(retmsg))
	case ERR_NONE:
		//send add success msg to client
		retmsg := new(MsgRetUserToBlack)
		retmsg.MsgId = MsgId_RetUserToBlack
		retmsg.Result = uint16(0)
		retmsg.Fuid = fuid
		client.send(Bytes(retmsg))
	}
}

func OnReqMoveFriendToGroup(client *Client, data []byte) {
	fuid := Uint64(data)
	group := data[8:]
	ret := gDataManager.moveFriendToGroup(client.uid, fuid, string(group))

	retmsg := new(MsgRetUserToBlack)
	retmsg.MsgId = MsgId_RetMoveFriendToGroup
	retmsg.Result = uint16(ret)
	retmsg.Fuid = fuid
	client.send(Bytes(retmsg))
}

func OnReqSetFriendVerifyType(client *Client, data []byte) {
	typ := data[0]
	ret := gDataManager.setFriendVerifyType(client.uid, typ)

	retmsg := new(MsgRetSetFriendVerifyType)
	retmsg.MsgId = MsgId_RetSetFriendVerifyType
	retmsg.Result = uint16(ret)
	client.send(Bytes(retmsg))
}
