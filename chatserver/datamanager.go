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

type dataManager interface {
	initialize()
	checkLogin(uid uint64, password string) bool
	//server op
	registerServer(addr string) bool
	incrServerClientCountBy(addr string, count int)
	getServerList() []string
	getServerCount() int
	setServerTTL(addr string, seconds int)
	checkServerTTL() int
	voteServerDie() int

	pullMsg(addr string, timeout int) []byte

	//user op
	setUserOnline()
	setUserOffline()
	isUserOnline()
	isUserExist()
	setUserState()

	//friend op
	reqAddFriend()
	agreeFriendReq()
	addFriend()
	deleteFriend()
	addFriendGroup()
	deleteFriendGroup()
	moveFriendToGroup()
	banFriend()
	unBanFriend()
	isFriend()
	setFriendVerify()
	getFriendVerify()
	setFriendAddSetting()
	getFriendAddSetting()

	//message op
	sendMsgToUser()
	sendMsgToRoom()

	//room op
	createRoom()
	deleteRoom()
	getRoomType()
	getRoomPassword()
	setRoomPassword()
	isRoomExist()
	isUserInRoom()
	addUserToRoom()
	removeUserToRoom()
	banUserInRoom()
	unBanUserInRoom()
	setRoomDescription()
	getRoomDescription()
	setRoomVerify()
	getRoomVerify()
	setRoomJoinSetting()
	getRoomJoinSetting()
}
