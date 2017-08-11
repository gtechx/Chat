package main

func OnReqFriendList(client *Client, data []byte) {
	flist, err := gDataManager.getFriendList(client.uid)

	if err != nil {
		reterr := NewErrorMsg(ERR_REDIS, MsgId_ReqFriendList)
		client.send(Byte(reterr))
		return
	}

	retmsg := new(MsgRetFriendList)
	retmsg.MsgId = MsgId_RetFriendList
	retmsg.GroupCount = byte(len(flist))

	data := []byte{}

	for groupname, uidarr := range flist {
		data = append(data, byte(len(groupname)))
		data = append(data, []byte(groupname)...)
		data = append(data, Bytes(uint16(len(uidarr))))

		for _, uid := range uidarr {
			data = append(data, Bytes(uid))
		}
	}

	retmsg.Data = data

	client.send(Bytes(retmsg))
}

func OnReqFriendAdd(client *Client, data []byte) {
	fuid := Uint64(data)
	group := data[8:]
	ret := gDataManager.addFriend(client.uid, fuid, group)

	switch ret {
	case ERR_FRIEND_IN_BLACKLIST:
		//do nothing
	case ERR_REDIS:
		reterr := NewErrorMsg(ERR_REDIS, MsgId_ReqFriendAdd)
		client.send(Byte(reterr))
	case ERR_FRIEND_EXIST:
		reterr := NewErrorMsg(ERR_FRIEND_EXIST, MsgId_ReqFriendAdd)
		client.send(Byte(reterr))
	case ERR_USER_NOT_EXIST:
		reterr := NewErrorMsg(ERR_USER_NOT_EXIST, MsgId_ReqFriendAdd)
		client.send(Byte(reterr))
	case ERR_FRIEND_MAX:
		reterr := NewErrorMsg(ERR_FRIEND_EXIST, MsgId_ReqFriendAdd)
		client.send(Byte(reterr))
	case ERR_FRIEND_ADD_NEED_REQ:
	case ERR_NONE:
	}
}

func OnReqFriendDel(client *Client, data []byte) {

}

func OnReqUserToBlack(client *Client, data []byte) {

}

func OnReqMoveFriendToGroup(client *Client, data []byte) {

}

func OnReqSetFriendVerifyType(client *Client, data []byte) {

}
