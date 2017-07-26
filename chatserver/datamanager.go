package main

import (
	"github.com/garyburd/redigo/redis"
	. "github.com/nature19862001/base/common"
)

type dataManager interface {
	//server op
	registerServer()
	getServerList()
	getServerCount()
	setServerTTL()
	checkServerTTL()
	voteServerDie()

	//user op
	setUserOnline()
	setUserOffline()
	isUserOnline()
	isUserExist()
	setUserState()

	//friend op
	addFriend()
	deleteFriend()
	addFriendGroup()
	deleteFriendGroup()
	moveFriendToGroup()
	banFriend()
	unBanFriend()
	isFriend()

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
}

var account map[uint64]string
var redisPool *redis.Pool

func init() {
	account = make(map[uint64]string)
	for i := 1000; i < 10000; i++ {
		account[uint64(i)] = Md5("1")
	}

	redisPool = &redis.Pool{
		MaxIdle:      3,
		IdleTimeout:  240 * time.Second,
		Dial:         redisDial,
		TestOnBorrow: redisOnBorrow,
	}
}

func redisDial() (redis.Conn, error) {
	c, err := redis.Dial("tcp", server)
	if err != nil {
		return nil, err
	}
	if _, err := c.Do("AUTH", password); err != nil {
		c.Close()
		return nil, err
	}
	if _, err := c.Do("SELECT", db); err != nil {
		c.Close()
		return nil, err
	}
	return c, nil
}

func redisOnBorrow(c redis.Conn, t time.Time) error {
	if time.Since(t) < time.Minute {
		return nil
	}
	_, err := c.Do("PING")
	return err
}

func checkLogin(uid uint64, password []byte) bool {
	pass, ok := account[uid]

	if ok && pass == string(password) {
		return true
	}

	return false
}

//server op
func getServerList() {

}

func getServerCount() {

}

func setServerTTL() {

}

func checkServerTTL() {

}

func voteServerDie() {

}

//user op
func setUserOnline() {

}

func setUserOffline() {

}

func isUserOnline() {

}

func isUserExist() {

}

func isFriend() {

}

func setUserState() {

}

func sendMsgToUser() {

}

func addFriend() {

}

func deleteFriend() {

}

func addFriendGroup() {

}

func deleteFriendGroup() {

}

func moveFriendToGroup() {

}

func banFriend() {

}

func unBanFriend() {

}

//room op

func sendMsgToRoom() {

}

func createRoom() {

}

func deleteRoom() {

}

func isRoomExist() {

}

func isUserInRoom() {

}

func addUserToRoom() {

}

func removeUserToRoom() {

}

func banUserInRoom() {

}

func unBanUserInRoom() {

}
