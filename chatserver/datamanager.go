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
	initialize() (bool, uint64)
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
	createUser(nickname, password, regip string) bool
	setUserOnline(uid uint64) bool
	setUserOffline(uid uint64)
	isUserOnline(uid uint64) bool
	isUserExist()
	setUserState()

	//friend op
	//reqAddFriend()
	//agreeFriendReq()
	addFriend()
	deleteFriend()
	getFriendList()
	addFriendGroup()
	deleteFriendGroup()
	moveFriendToGroup()
	banFriend()
	unBanFriend()
	isFriend()
	setFriendVerify()
	getFriendVerify()
	setFriendVerifyType()
	getFriendVerifyType()

	//message op
	sendMsgToUser()
	sendMsgToRoom()

	//room op
	createRoom()
	deleteRoom()
	getRoomList()
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
	setRoomVerifyType()
	getRoomVerifyType()
}
