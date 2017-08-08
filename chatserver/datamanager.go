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
	initialize() bool
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
	createUser(nickname, password, regip string) (bool, uint64)
	setUserOnline(uid uint64) bool
	setUserOffline(uid uint64)
	isUserOnline(uid uint64) bool
	isUserExist(uid uint64) bool
	setUserState()
	addUserToBlacklist(uid, uuid uint64) int
	isUserInBlacklist(uid, uuid uint64) bool

	//friend op
	//reqAddFriend()
	//agreeFriendReq()
	addFriendReq() int
	addFriend(uid, fuid uint64, group string) int
	deleteFriend(uid, fuid uint64) int
	getFriendList() map[string][]uint64
	addFriendGroup(uid uint64, groupname string) int
	deleteFriendGroup(uid uint64, groupname string) int
	moveFriendToGroup(uid, fuid uint64, srcgroup, destgroup string) int
	banFriend()
	unBanFriend()
	isFriend(uid, fuid uint64) bool
	setFriendVerifyType(uid uint64, vtype byte) int
	getFriendVerifyType() byte

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
