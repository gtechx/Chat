package main

import (
	//"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	. "github.com/nature19862001/base/common"
	"strings"
	"time"
)

var serverListKeyName string = "serverlist"
var userOnlineKeyName string = "user:online"
var defaultGroupName string = "我的好友"

//key							field		field	...
//uid				hashes		nickname	password
//fgroup:uid		sets
//friend:uid		hashes		fuid
//friend:group:uid	hashes		group:(n)	groupname
//black:uid			sets
//freq				hashes		uid:fuid

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

	ret, err := conn.Do("BLPOP", "msg:"+addr, timeout)

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

	ret, err = conn.Do("HMSET", uid, "nickname", nickname, "password", password, "regip", regip, "regdate", time.Now().Unix(), "maxfriends", 1000, "headurl", "", "desc", "")

	if err != nil {
		fmt.Println("createUser error:", err.Error())
		return false, 0
	}

	ret, err = conn.Do("SADD", "fgroup:"+String(uid), defaultGroupName)

	if err != nil {
		fmt.Println("createUser error:", err.Error())
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

	ret, err := conn.Do("EXISTS", uid)

	if err != nil {
		fmt.Println("isUserExist error:", err.Error())
		return false
	}

	return Bool(ret)
}

func (this *redisDataManager) setUserState() {

}

//friend op
// func (this *redisDataManager) reqAddFriend() {

// }

// func (this *redisDataManager) agreeFriendReq() {

// }

func (this *redisDataManager) addUserToBlacklist(uid, uuid uint64) int {
	// if !isUserExist(uuid) {
	// 	return ERR_USER_NOT_EXIST
	// }
	conn := this.redisPool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("HDEL", "friend:"+String(uid), uuid)
	//conn.Send("HDEL", "fgroup:"+String(uid)+":"+group, fuid)
	conn.Send("SADD", "black:"+String(uuid), uid)
	_, err := conn.Do("EXEC")

	if err != nil {
		fmt.Println("addUserToBlacklist error:", err.Error())
		return ERR_REDIS
	}

	return ERR_NONE
}

func (this *redisDataManager) isUserInBlacklist(uid, uuid uint64) bool {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("SISMEMBER", "black:"+String(uid), uuid)

	if err != nil {
		fmt.Println("isUserInBlacklist error:", err.Error())
		return false
	}

	return Bool(ret)
}

//check if uid2 is in uid1's blacklist
func (this *redisDataManager) isInBlacklist(uid1, uid2 uint64, conn redis.Conn) bool {
	ret, err := conn.Do("SISMEMBER", "black:"+String(uid1), uid2)

	if err != nil {
		fmt.Println("isInBlacklist error:", err.Error())
		return true
	}

	return Bool(ret)
}

func (this *redisDataManager) addFriendReq(uid, fuid uint64, group string) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	//check if in black list
	// ret, err := conn.Do("SISMEMBER", "black:"+String(fuid), uid)

	// if err != nil {
	// 	fmt.Println("addFriend error:", err.Error())
	// 	return ERR_REDIS
	// }

	if this.isInBlacklist(fuid, uid, conn) {
		return ERR_FRIEND_IN_BLACKLIST
	}

	//check if group exists
	_, err := conn.Do("HSET", "freq", String(uid)+":"+String(fuid), group)

	if err != nil {
		fmt.Println("addFriendReq error:", err.Error())
		return -1
	}

	return ERR_NONE
}

func (this *redisDataManager) addFriend(uid, fuid uint64, group string) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	if this.isInBlacklist(fuid, uid, conn) {
		return ERR_FRIEND_IN_BLACKLIST
	}

	//check if friend is exist
	ret, err := conn.Do("SISMEMBER", "friend:"+String(uid), fuid)

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return ERR_REDIS
	}

	if Bool(ret) {
		return ERR_FRIEND_EXIST
	}

	//check if fuid is exists
	ret, err = conn.Do("HEXISTS", fuid)

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return ERR_REDIS
	}

	if Bool(ret) {
		return ERR_USER_NOT_EXIST
	}

	//check if friend count is max
	ret, err = conn.Do("HGET", uid, "maxfriends")

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return ERR_REDIS
	}

	maxfriends := Int(ret)

	ret, err = conn.Do("SCARD", "friend:"+String(uid))

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return ERR_REDIS
	}

	curfriendcount := Int(ret)

	if curfriendcount >= maxfriends {
		return ERR_FRIEND_MAX
	}

	//check if group exists
	ret, err = conn.Do("SISMEMBER", "fgroup:"+String(uid), group)

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return ERR_REDIS
	}

	if Bool(ret) != true {
		//return ERR_FRIEND_GROUP_NOT_EXIST
		group = defaultGroupName
	}

	ret, err = conn.Do("HGET", "freq", String(fuid)+":"+String(uid))

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return ERR_REDIS
	}

	if ret == nil {
		//need friend request to fuid
		return ERR_FRIEND_ADD_NEED_REQ
	}

	fuidgroup := String(ret)

	ret, err = conn.Do("HGET", "friend:group:"+String(uid), group)

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return ERR_REDIS
	}

	groupuidcount := 0

	if ret != nil {
		groupuidcount = Int(ret)
	}

	ret, err = conn.Do("HGET", "friend:group:"+String(fuid), fuidgroup)

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return ERR_REDIS
	}

	fuidgroupuidcount := 0

	if ret != nil {
		fuidgroupuidcount = Int(ret)
	}

	conn.Send("MULTI")
	conn.Send("HSET", "friend:"+String(uid), fuid, group+":"+String(groupuidcount))
	conn.Send("HSET", "friend:group:"+String(uid), group+":"+String(groupuidcount), fuid)
	conn.Send("HINCRBY", "friend:group:"+String(uid), group, 1)

	conn.Send("HSET", "friend:"+String(fuid), uid, fuidgroup+":"+String(fuidgroupuidcount))
	conn.Send("HSET", "friend:group:"+String(fuid), fuidgroup+":"+String(fuidgroupuidcount), uid)
	conn.Send("HINCRBY", "friend:group:"+String(fuid), fuidgroup, 1)

	conn.Send("HDEL", "freq", String(fuid)+":"+String(uid))
	_, err = conn.Do("EXEC")

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return ERR_REDIS
	}

	return ERR_NONE
}

func (this *redisDataManager) deleteFriend(uid, fuid uint64) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	//get group
	ret, err := conn.Do("HGET", "friend:"+String(uid), fuid)

	if err != nil {
		fmt.Println("deleteFriend error:", err.Error())
		return ERR_REDIS
	}

	if ret == nil {
		return ERR_FRIEND_NOT_EXIST
	}

	groupstr := String(ret)
	groupstrarr := strings.Split(groupstr, ":")
	group := groupstrarr[0]

	conn.Send("MULTI")
	conn.Send("HDEL", "friend:"+String(uid), fuid)
	conn.Send("HDEL", "friend:group:"+String(uid), groupstr)
	conn.Send("HINCRBY", "friend:group:"+String(uid), group, -1)
	//conn.Send("SREM", "fgroup:"+String(uid)+":"+group, fuid)
	_, err = conn.Do("EXEC")

	if err != nil {
		fmt.Println("deleteFriend error:", err.Error())
		return ERR_REDIS
	}

	return ERR_NONE
}

func (this *redisDataManager) getFriendList(uid uint64) (map[string][]uint64, int) {
	conn := this.redisPool.Get()
	defer conn.Close()

	//
	ret, err := conn.Do("HGETALL", "friend:group:"+String(uid))

	if err != nil {
		fmt.Println("getFriendList error:", err.Error())
		return map[string][]uint64{}, ERR_REDIS
	}

	if ret == nil {
		return map[string][]uint64{}, ERR_NONE
	}

	dataarr, err := redis.Strings(ret, err)
	retdata := map[string][]uint64{}

	for i := 0; i < len(dataarr); i += 2 {
		groupstr := dataarr[i]
		groupstrarr := strings.Split(groupstr, ":")
		if len(groupstrarr) <= 1 {
			continue
		}
		groupname := groupstrarr[0]

		_, ok := retdata[groupname]

		fuid := Uint64(groupstrarr[1])
		if !ok {
			retdata[groupname] = make([]uint64, 0)
		}
		retdata[groupname] = append(retdata[groupname], fuid)
	}

	return retdata, ERR_NONE
}

func (this *redisDataManager) addFriendGroup(uid uint64, groupname string) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("SADD", "fgroup:"+String(uid), groupname)

	if err != nil {
		fmt.Println("addFriendGroup error:", err.Error())
		return ERR_REDIS
	}

	if Int(ret) != 1 {
		return ERR_FRIEND_GROUP_EXIST
	}

	return ERR_NONE
}

func (this *redisDataManager) deleteFriendGroup(uid uint64, groupname string) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("HGET", "friend:group:"+String(uid), groupname)

	if err != nil {
		fmt.Println("deleteFriendGroup error:", err.Error())
		return ERR_REDIS
	}

	if ret != nil && Int(ret) > 0 {
		return ERR_FRIEND_GROUP_USER_NOT_EMPTY
	}

	ret, err = conn.Do("SREM", "fgroup:"+String(uid), groupname)

	if err != nil {
		fmt.Println("deleteFriendGroup error:", err.Error())
		return ERR_REDIS
	}

	// if Int(ret) != 1 {
	// 	return 0
	// }

	return ERR_NONE
}

func (this *redisDataManager) moveFriendToGroup(uid, fuid uint64, destgroup string) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	//check if friend is exist
	ret, err := conn.Do("HEXISTS", "friend:"+String(uid), fuid)

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return ERR_REDIS
	}

	if !Bool(ret) {
		return ERR_FRIEND_NOT_EXIST
	}

	//check if group exists
	// ret, err = conn.Do("SISMEMBER", "fgroup:"+String(uid), srcgroup)

	// if err != nil {
	// 	fmt.Println("moveFriendToGroup error:", err.Error())
	// 	return ERR_REDIS
	// }

	// if Bool(ret) != true {
	// 	return ERR_FRIEND_GROUP_NOT_EXIST
	// }

	ret, err = conn.Do("SISMEMBER", "fgroup:"+String(uid), destgroup)

	if err != nil {
		fmt.Println("moveFriendToGroup error:", err.Error())
		return ERR_REDIS
	}

	if Bool(ret) != true {
		return ERR_FRIEND_GROUP_NOT_EXIST
	}

	//get group fuid current in
	ret, err = conn.Do("HGET", "friend:"+String(uid), fuid)

	if err != nil {
		fmt.Println("moveFriendToGroup error:", err.Error())
		return ERR_REDIS
	}

	groupstr := String(ret)
	groupstrarr := strings.Split(groupstr, ":")
	curgroup := groupstrarr[0]

	if curgroup == destgroup {
		//if already in destgroup, return
		return ERR_NONE
	}

	//get destgroup count
	ret, err = conn.Do("HGET", "friend:group:"+String(uid), destgroup)

	if err != nil {
		fmt.Println("moveFriendToGroup error:", err.Error())
		return ERR_REDIS
	}

	destgroupcount := 0

	if ret != nil {
		destgroupcount = Int(ret)
	}

	conn.Send("MULTI")
	//remove from current group
	conn.Send("HDEL", "friend:"+String(uid), fuid)
	conn.Send("HDEL", "friend:group:"+String(uid), groupstr)
	conn.Send("HINCRBY", "friend:group:"+String(uid), curgroup, -1)

	//add to dest group
	conn.Send("HSET", "friend:"+String(uid), fuid, destgroup+":"+String(destgroupcount))
	conn.Send("HSET", "friend:group:"+String(uid), destgroup+":"+String(destgroupcount), fuid)
	conn.Send("HINCRBY", "friend:group:"+String(uid), destgroup, 1)

	// conn.Send("SADD", "fgroup:"+String(uid)+":"+destgroup, fuid)
	// conn.Send("SREM", "fgroup:"+String(uid)+":"+srcgroup, fuid)
	_, err = conn.Do("EXEC")

	if err != nil {
		fmt.Println("deleteFriend error:", err.Error())
		return ERR_REDIS
	}

	// ret, err = conn.Do("SADD", "fgroup:"+String(uid)+":"+destgroup, fuid)

	// if err != nil {
	// 	fmt.Println("moveFriendToGroup error:", err.Error())
	// 	return -1
	// }

	// if Int(ret) != 1 {
	// 	return 0
	// }

	// ret, err = conn.Do("SREM", "fgroup:"+String(uid)+":"+srcgroup, fuid)

	// if err != nil {
	// 	fmt.Println("moveFriendToGroup error:", err.Error())
	// 	return -1
	// }

	// if Int(ret) != 1 {
	// 	return 0
	// }

	return ERR_NONE
}

func (this *redisDataManager) banFriend(uid, fuid uint64) {

}

func (this *redisDataManager) unBanFriend(uid, fuid uint64) {

}

func (this *redisDataManager) isFriend(uid, fuid uint64) bool {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("HEXISTS", "friend:"+String(uid), fuid)

	if err != nil {
		fmt.Println("getFriendVerifyType error:", err.Error())
		return false
	}

	return Bool(ret)
}

func (this *redisDataManager) setFriendVerifyType(uid uint64, vtype byte) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	_, err := conn.Do("HSET", uid, "vtype", vtype)

	if err != nil {
		fmt.Println("setFriendVerifyType error:", err.Error())
		return ERR_REDIS
	}

	return ERR_NONE
}

func (this *redisDataManager) getFriendVerifyType(uid uint64) byte {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("HGET", uid, "vtype")

	if err != nil {
		fmt.Println("getFriendVerifyType error:", err.Error())
		return VERIFY_TYPE_ERR
	}

	return Byte(ret)
}

//message op
func (this *redisDataManager) sendMsgToUser(uid uint64, data []byte) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("HGET", uid, "online")

	if err != nil {
		fmt.Println("sendMsgToUser error:", err.Error())
		return ERR_REDIS
	}

	if ret != nil {
		// if online
		serveraddr := String(ret)
		ret, err = conn.Do("RPUSH", "msg:"+serveraddr, data)

		if err != nil {
			fmt.Println("sendMsgToUser error:", err.Error())
			return ERR_REDIS
		}
	} else {
		//else not online
		ret, err = conn.Do("RPUSH", "offline:"+String(uid), data)

		if err != nil {
			fmt.Println("sendMsgToUser error:", err.Error())
			return ERR_REDIS
		}
	}

	return ERR_NONE
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
