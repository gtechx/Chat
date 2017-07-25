package main

import (
	"github.com/garyburd/redigo/redis"
	. "github.com/nature19862001/base/common"
)

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

func setUserOnline() {

}

func setUserOffline() {

}

func setUserState() {

}

func sendMsgToUser() {

}

func sendMsgToRoom() {

}

func createRoom() {

}

func deleteRoom() {

}

func addFriend() {

}

func deleteFriend() {

}
