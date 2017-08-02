package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	. "github.com/nature19862001/base/common"
	"time"
)

var serverListKeyName string = "serverlist"
var userOnlineKeyName string = "user:online"

type redisDataManager struct {
	redisPool *redis.Pool
}

func (this *redisDataManager) initialize() bool {
	this.redisPool = &redis.Pool{
		MaxIdle:      3,
		IdleTimeout:  240 * time.Second,
		Dial:         redisDial,
		TestOnBorrow: redisOnBorrow,
	}

	conn := this.redisPool.Get()
	defer conn.Close()
	n, err := conn.Do("EXISTS", "UID")

	if err != nil {
		fmt.Println("redis server error:", err.Error())
		return false
	}

	if !Bool(n) {
		_, err = conn.Do("SET", "UID", 10000)

		if err != nil {
			fmt.Println("redis server error:", err.Error())
			return false
		}
	}

	return true
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
	_, err := conn.Do("ZADD", serverListKeyName, 0, addr)

	if err != nil {
		fmt.Println("register server error:", err.Error())
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
func (this *redisDataManager) createUser(nickname, password, regip string) (bool, uint64) {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("INCR", "UID")

	if err != nil {
		fmt.Println("createUser error:", err.Error())
		return false, 0
	}

	uid := Uint64(ret)

	ret, err = conn.Do("HMSET", uid, "nickname", nickname, "password", password, "regip", regip, "regdate", time.Now().Unix())

	if err != nil {
		fmt.Println("createUser error:", err.Error())
		return false, 0
	}

	return true, uid
}

func (this *redisDataManager) setUserOnline(uid uint64) bool {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("HSET", uid, "online", serverAddr)

	if err != nil {
		fmt.Println("setUserOnline error:", err.Error())
		return false
	}

	if ret == nil {
		return false
	}

	return true
}

func (this *redisDataManager) setUserOffline(uid uint64) {
	conn := this.redisPool.Get()
	defer conn.Close()

	_, err := conn.Do("HDEL", uid, "online")

	if err != nil {
		fmt.Println("setUserOffline error:", err.Error())
	}
}

func (this *redisDataManager) isUserOnline(uid uint64) bool {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("HEXISTS", uid, "online")

	if err != nil {
		fmt.Println("isUserOnline error:", err.Error())
		return false
	}

	return Int(ret) == 1
}

func (this *redisDataManager) isUserExist(uid uint64) bool {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("TYPE", uid)

	if err != nil {
		fmt.Println("isUserExist error:", err.Error())
		return false
	}

	return String(ret) == "hash"
}

func (this *redisDataManager) setUserState() {

}

//friend op
// func (this *redisDataManager) reqAddFriend() {

// }

// func (this *redisDataManager) agreeFriendReq() {

// }

func (this *redisDataManager) addFriendReq(uid, fuid uint64, group string) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	//check if group exists
	ret, err := conn.Do("HSET", "freq", String(uid)+":"+String(fuid), group)

	if err != nil {
		fmt.Println("addFriendReq error:", err.Error())
		return -1
	}

	return 1
}

func (this *redisDataManager) addFriend(uid, fuid uint64, group string) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	//check if group exists
	ret, err := conn.Do("SISMEMBER", "fgroup:"+String(uid), group)

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return -1
	}

	if Bool(ret) != true {
		return 2
	}

	ret, err = conn.Do("HEXISTS", "freq", String(uid)+":"+String(fuid))

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return -1
	}

	if !Bool(ret) {
		//need friend request to fuid
		return 3
	}

	//add friend
	ret, err = conn.Do("SADD", "friend:"+String(uid), fuid)

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return -1
	}

	if Int(ret) != 1 {
		return 0
	}

	ret, err = conn.Do("SADD", "fgroup:"+String(uid)+":"+group, fuid)

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return -1
	}

	if Int(ret) != 1 {
		return 0
	}

	//add inverse friend
	ret, err = conn.Do("SADD", "friend:"+String(fuid), uid)

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return -1
	}

	if Int(ret) != 1 {
		return 0
	}

	ret, err = conn.Do("SADD", "fgroup:"+String(fuid)+":"+group, uid)

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return -1
	}

	if Int(ret) != 1 {
		return 0
	}

	return 1
}

func (this *redisDataManager) deleteFriend(uid, fuid uint64) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("SREM", "friend:"+String(uid), fuid)

	if err != nil {
		fmt.Println("deleteFriend error:", err.Error())
		return -1
	}

	if Int(ret) != 1 {
		return 0
	}

	ret, err = conn.Do("SREM", "fgroup:"+String(uid)+":"+group, fuid)

	if err != nil {
		fmt.Println("deleteFriend error:", err.Error())
		return -1
	}

	if Int(ret) != 1 {
		return 0
	}

	return 1
}

func (this *redisDataManager) getFriendList() map[string][]uint64 {
	//
}

func (this *redisDataManager) addFriendGroup(uid uint64, groupname string) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("SADD", "fgroup:"+String(uid), groupname)

	if err != nil {
		fmt.Println("addFriendGroup error:", err.Error())
		return -1
	}

	if Int(ret) != 1 {
		return 0
	}

	return 1
}

func (this *redisDataManager) deleteFriendGroup(uid uint64, groupname string) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("SREM", "fgroup:"+String(uid), groupname)

	if err != nil {
		fmt.Println("deleteFriendGroup error:", err.Error())
		return -1
	}

	if Int(ret) != 1 {
		return 0
	}

	return 1
}

func (this *redisDataManager) moveFriendToGroup(uid, fuid uint64, srcgroup, destgroup string) int {
	//check if group exists
	ret, err := conn.Do("SISMEMBER", "fgroup:"+String(uid), srcgroup)

	if err != nil {
		fmt.Println("moveFriendToGroup error:", err.Error())
		return -1
	}

	if Bool(ret) != true {
		return 2
	}

	ret, err = conn.Do("SISMEMBER", "fgroup:"+String(uid), destgroup)

	if err != nil {
		fmt.Println("moveFriendToGroup error:", err.Error())
		return -1
	}

	if Bool(ret) != true {
		return 2
	}

	ret, err = conn.Do("SADD", "fgroup:"+String(uid)+":"+destgroup, fuid)

	if err != nil {
		fmt.Println("moveFriendToGroup error:", err.Error())
		return -1
	}

	if Int(ret) != 1 {
		return 0
	}

	ret, err = conn.Do("SREM", "fgroup:"+String(uid)+":"+srcgroup, fuid)

	if err != nil {
		fmt.Println("moveFriendToGroup error:", err.Error())
		return -1
	}

	if Int(ret) != 1 {
		return 0
	}

	return 1
}

func (this *redisDataManager) banFriend(uid, fuid uint64) {

}

func (this *redisDataManager) unBanFriend(uid, fuid uint64) {

}

func (this *redisDataManager) isFriend(uid, fuid uint64) {

}

func (this *redisDataManager) setFriendVerify() {

}

func (this *redisDataManager) getFriendVerify() {

}

func (this *redisDataManager) setFriendVerifyType() {

}

func (this *redisDataManager) getFriendVerifyType() {

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

func (this *redisDataManager) getRoomList() {

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

func (this *redisDataManager) setRoomVerifyType() {

}

func (this *redisDataManager) getRoomVerifyType() {

}
