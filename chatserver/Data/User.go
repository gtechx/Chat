package data

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	. "github.com/nature19862001/base/common"
)

func (this *RedisDataManager) IsAccountExist(account string) (bool, error) {
	conn := this.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("EXISTS", uid)
	return Bool(ret), err
}

func (this *RedisDataManager) CreateAccount(account, password, regip string) error {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("INCR", "UID")

	if err != nil {
		return err
	}

	uid := Uint64(ret)

	conn.Send("MULTI")
	conn.Send("HMSET", uid, "account", account, "password", password, "regip", regip, "regdate", time.Now().Unix())
	conn.Send("SADD", "fgroup:"+String(uid), defaultGroupName)
	conn.Send("SADD", "user", uid)

	_, err = conn.Do("EXEC")
	return err
}

func (this *RedisDataManager) CreateAccount(account, password, regip string) error {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("INCR", "UID")

	if err != nil {
		return err
	}

	uid := Uint64(ret)

	conn.Send("MULTI")
	conn.Send("HMSET", uid, "nickname", nickname, "password", password, "regip", regip, "regdate", time.Now().Unix(), "maxfriends", 1000, "headurl", "", "desc", "")
	conn.Send("SADD", "fgroup:"+String(uid), defaultGroupName)
	conn.Send("SADD", "user", uid)

	_, err = conn.Do("EXEC")

	if err != nil {
		fmt.Println("createUser error:", err.Error())
		return false, 0
	}

	return true, uid
}

func (this *RedisDataManager) checkLogin(uid uint64, password string) int {
	if password == "" {
		return ERR_PASSWORD_INVALID
	}
	conn := this.redisPool.Get()
	defer conn.Close()
	fmt.Println("checklogin:", uid, password)
	ret, err := conn.Do("HGET", uid, "password")

	if err != nil {
		fmt.Println("checkLogin error:", err.Error())
		return ERR_REDIS
	}

	if ret == nil {
		return ERR_USER_NOT_EXIST
	}

	if String(ret) != password {
		return ERR_PASSWORD_INVALID
	}

	return ERR_NONE
}

func (this *RedisDataManager) setUserOnline(uid uint64) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("HSET", uid, "online", serverAddr)
	conn.Send("SADD", "online", uid)

	_, err := conn.Do("EXEC")

	if err != nil {
		fmt.Println("setUserOnline error:", err.Error())
		return ERR_REDIS
	}

	// ret, err := conn.Do("HSET", uid, "online", serverAddr)

	// if err != nil {
	// 	fmt.Println("setUserOnline error:", err.Error())
	// 	return false
	// }

	// if ret == nil {
	// 	return false
	// }

	return ERR_NONE
}

func (this *RedisDataManager) isUserOnline(uid uint64) bool {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("HEXISTS", uid, "online")

	if err != nil {
		fmt.Println("isUserOnline error:", err.Error())
		return false
	}

	return Int(ret) == 1
}

func (this *RedisDataManager) isUserExist(uid uint64) bool {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("EXISTS", uid)

	if err != nil {
		fmt.Println("isUserExist error:", err.Error())
		return false
	}

	return Bool(ret)
}

func (this *RedisDataManager) setUserState() {

}

func (this *RedisDataManager) addUserToBlacklist(uid, uuid uint64) int {
	// if !isUserExist(uuid) {
	// 	return ERR_USER_NOT_EXIST
	// }
	conn := this.redisPool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	//conn.Send("HDEL", "friend:"+String(uid), uuid)
	//conn.Send("HDEL", "fgroup:"+String(uid)+":"+group, fuid)
	conn.Send("SADD", "black:"+String(uid), uuid)
	_, err := conn.Do("EXEC")

	if err != nil {
		fmt.Println("addUserToBlacklist error:", err.Error())
		return ERR_REDIS
	}

	return ERR_NONE
}

func (this *RedisDataManager) removeUserInBlacklist(uid, uuid uint64) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	//conn.Send("HDEL", "friend:"+String(uid), uuid)
	//conn.Send("HDEL", "fgroup:"+String(uid)+":"+group, fuid)
	conn.Send("SREM", "black:"+String(uid), uuid)
	_, err := conn.Do("EXEC")

	if err != nil {
		fmt.Println("removeUserInBlacklist error:", err.Error())
		return ERR_REDIS
	}

	return ERR_NONE
}

func (this *RedisDataManager) isUserInBlacklist(uid, uuid uint64) bool {
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
func (this *RedisDataManager) isInBlacklist(uid1, uid2 uint64, conn redis.Conn) bool {
	ret, err := conn.Do("SISMEMBER", "black:"+String(uid1), uid2)

	if err != nil {
		fmt.Println("isInBlacklist error:", err.Error())
		return true
	}

	return Bool(ret)
}
