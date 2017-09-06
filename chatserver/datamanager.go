package main

import (
//. "github.com/nature19862001/base/common"
)

// var account map[uint64]string

// func init() {
// 	account = make(map[uint64]string)
// 	for i := 1000; i < 10000; i++ {
// 		account[uint64(i)] = Md5("1")
// 	}
// }

type dataManager interface {
	initialize() bool
	checkLogin(uid uint64, password string) int
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
	addAdmin(uid, uuid uint64, privilege uint32) int
	removeAdmin(uid, uuid uint64) int

	getAdminList(uid uint64) ([]uint64, int)
	getUserList(uid uint64) ([]uint64, int)
	getUserOnline(uid uint64) ([]uint64, int)

	createUser(nickname, password, regip string) (bool, uint64)
	setUserOnline(uid uint64) int
	setUserOffline(uid uint64)

	isUserOnline(uid uint64) bool
	isUserExist(uid uint64) bool
	setUserState()
	addUserToBlacklist(uid, uuid uint64) int
	removeUserInBlacklist(uid, uuid uint64) int
	isUserInBlacklist(uid, uuid uint64) bool

	//friend op
	//reqAddFriend()
	//agreeFriendReq()
	addFriendReq(uid, fuid uint64, group string) int
	addFriend(uid, fuid uint64, group string) int
	deleteFriend(uid, fuid uint64) int
	getFriendList(uid uint64) (map[string][]uint64, int)
	addFriendGroup(uid uint64, groupname string) int
	deleteFriendGroup(uid uint64, groupname string) int
	getGroupOfFriend(uid, fuid uint64) (string, int)
	moveFriendToGroup(uid, fuid uint64, destgroup string) int
	banFriend(uid, fuid uint64)
	unBanFriend(uid, fuid uint64)
	isFriend(uid, fuid uint64) bool
	setFriendVerifyType(uid uint64, vtype byte) int
	getFriendVerifyType(uid uint64) byte

	//message op
	sendMsgToUser(uid uint64, data []byte) int
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

	//app op
	createApp(uid uint64, name, password, desc, iconurl string) int
	deleteApp(uid uint64, name string) int
	setAppOnline(appname string) int
	setAppOffline(appname string)
	isAppUser(appname string, uid uint64) bool
	createAppUser(puid uint64, appname, nickname, password, regip string) (int, uint64)

	//app
	checkAppLogin(appname, password string) int
	setAppVerifyData(uuid string, uid uint64) int
	verifyAppLoginData(uuid string, uid uint64) bool
}
