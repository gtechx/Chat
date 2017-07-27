package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	. "github.com/nature19862001/base/common"
	"time"
)

var serverListKeyName string = "serverlist"

type redisDataManager struct {
	redisPool *redis.Pool
}

func (this *redisDataManager) initialize() {
	this.redisPool = &redis.Pool{
		MaxIdle:      3,
		IdleTimeout:  240 * time.Second,
		Dial:         redisDial,
		TestOnBorrow: redisOnBorrow,
	}
}

func (this *redisDataManager) checkLogin(uid uint64, password string) bool {
	return true
}

func redisDial() (redis.Conn, error) {
	c, err := redis.Dial("tcp", redisAddr)
	if err != nil {
		return nil, err
	}
	if _, err := c.Do("AUTH", "ztgame@123"); err != nil {
		c.Close()
		return nil, err
	}
	// if _, err := c.Do("SELECT", db); err != nil {
	// 	c.Close()
	// 	return nil, err
	// }
	return c, nil
}

func redisOnBorrow(c redis.Conn, t time.Time) error {
	if time.Since(t) < time.Minute {
		return nil
	}
	_, err := c.Do("PING")
	return err
}

//server op
func (this *redisDataManager) registerServer(addr string) bool {
	conn := this.redisPool.Get()
	defer conn.Close()
	n, err := conn.Do("ZADD", serverListKeyName, 0, addr)

	if err != nil {
		fmt.Println("register server error:", err.Error())
		return false
	}

	if n == nil {
		return false
	}

	if Int(n) <= 0 {
		return false
	}

	return true
}

func (this *redisDataManager) incrServerClientCountBy(addr string, count int) {
	conn := this.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("ZINCRBY", serverListKeyName, count, addr)

	if err != nil {
		fmt.Println("incrServerClientCountBy error:", err.Error())
	}
}

func (this *redisDataManager) getServerList() []string {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("ZRANGE", serverListKeyName, 0, -1)

	if err != nil {
		fmt.Println("getServerList error:", err.Error())
		return []string{}
	}

	slist, _ := redis.Strings(ret, err)
	return slist
}

func (this *redisDataManager) getServerCount() int {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("ZCARD", serverListKeyName)

	if err != nil {
		fmt.Println("getServerCount error:", err.Error())
		return 0
	}

	return Int(ret)
}

func (this *redisDataManager) setServerTTL(addr string, seconds int) {
	conn := this.redisPool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", "ttl"+addr, 0, "EX", seconds)

	if err != nil {
		fmt.Println("setServerTTL error:", err.Error())
		return
	}
}

func (this *redisDataManager) checkServerTTL() int {
	return 1
}

func (this *redisDataManager) voteServerDie() int {
	return 0
}

func (this *redisDataManager) pullMsg(addr string, timeout int) []byte {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("BLPOP", "msg"+addr, timeout)

	if err != nil {
		fmt.Println("pullMsg error:", err.Error())
		return nil
	}

	if ret == nil {
		return nil
	} else {
		return Bytes(ret)
	}
}

//user op
func (this *redisDataManager) setUserOnline() {

}

func (this *redisDataManager) setUserOffline() {

}

func (this *redisDataManager) isUserOnline() {

}

func (this *redisDataManager) isUserExist() {

}

func (this *redisDataManager) setUserState() {

}

//friend op
func (this *redisDataManager) reqAddFriend() {

}

func (this *redisDataManager) addFriend() {

}

func (this *redisDataManager) agreeFriendReq() {

}

func (this *redisDataManager) deleteFriend() {

}

func (this *redisDataManager) addFriendGroup() {

}

func (this *redisDataManager) deleteFriendGroup() {

}

func (this *redisDataManager) moveFriendToGroup() {

}

func (this *redisDataManager) banFriend() {

}

func (this *redisDataManager) unBanFriend() {

}

func (this *redisDataManager) isFriend() {

}

func (this *redisDataManager) setFriendVerify() {

}

func (this *redisDataManager) getFriendVerify() {

}

func (this *redisDataManager) setFriendAddSetting() {

}

func (this *redisDataManager) getFriendAddSetting() {

}

//message op
func (this *redisDataManager) sendMsgToUser() {

}

func (this *redisDataManager) sendMsgToRoom() {

}

//room op
func (this *redisDataManager) createRoom() {

}

func (this *redisDataManager) deleteRoom() {

}

func (this *redisDataManager) getRoomType() {

}

func (this *redisDataManager) getRoomPassword() {

}

func (this *redisDataManager) setRoomPassword() {

}

func (this *redisDataManager) isRoomExist() {

}

func (this *redisDataManager) isUserInRoom() {

}

func (this *redisDataManager) addUserToRoom() {

}

func (this *redisDataManager) removeUserToRoom() {

}

func (this *redisDataManager) banUserInRoom() {

}

func (this *redisDataManager) unBanUserInRoom() {

}

func (this *redisDataManager) setRoomDescription() {

}

func (this *redisDataManager) getRoomDescription() {

}

func (this *redisDataManager) setRoomVerify() {

}

func (this *redisDataManager) getRoomVerify() {

}

func (this *redisDataManager) setRoomJoinSetting() {

}

func (this *redisDataManager) getRoomJoinSetting() {

}
