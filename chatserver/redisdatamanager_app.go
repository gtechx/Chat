package main

import (
	//"errors"
	"fmt"
	//"github.com/garyburd/redigo/redis"
	. "github.com/nature19862001/base/common"
	//"strings"
	"time"
)

//app op
func (this *redisDataManager) createApp(uid uint64, name, password, desc, iconurl string) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("SISMEMBER", "app", name)

	if err != nil {
		fmt.Println("createApp error:", err.Error())
		return ERR_REDIS
	}

	if Bool(ret) {
		return ERR_APP_EXIST
	}

	conn.Send("SADD", "app", name)
	conn.Send("HMSET", "app:"+name, "password", password, "desc", desc, "iconurl", iconurl, "regdate", time.Now().Unix(), "maxfriends", 1000)

	_, err = conn.Do("EXEC")

	if err != nil {
		fmt.Println("createApp error:", err.Error())
		return ERR_REDIS
	}

	return ERR_NONE
}

func (this *redisDataManager) deleteApp(uid uint64, name string) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("SISMEMBER", "app", name)

	if err != nil {
		fmt.Println("deleteApp error:", err.Error())
		return ERR_REDIS
	}

	if !Bool(ret) {
		return ERR_APP_NOT_EXIST
	}

	conn.Send("SREM", "app", name)
	conn.Send("DEL", "app:"+name)

	_, err = conn.Do("EXEC")

	if err != nil {
		fmt.Println("deleteApp error:", err.Error())
		return ERR_REDIS
	}

	return ERR_NONE
}

func (this *redisDataManager) setAppOnline(appname string) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("HSET", "app:"+appname, "online", serverAddr)
	conn.Send("SADD", "app:online", appname)

	_, err := conn.Do("EXEC")

	if err != nil {
		fmt.Println("setAppOnline error:", err.Error())
		return ERR_REDIS
	}
	return ERR_NONE
}

func (this *redisDataManager) setAppOffline(appname string) {
	conn := this.redisPool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("HDEL", "app:"+appname, "online")
	conn.Send("SREM", "app:online", appname)

	_, err := conn.Do("EXEC")

	if err != nil {
		fmt.Println("setAppOffline error:", err.Error())
	}
}

func (this *redisDataManager) createAppUser(puid uint64, appname, nickname, password, regip string) (int, uint64) {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("INCR", "UID")

	if err != nil {
		fmt.Println("createAppUser error:", err.Error())
		return ERR_REDIS, 0
	}

	uid := Uint64(ret)

	conn.Send("MULTI")
	conn.Send("HMSET", uid, "nickname", nickname, "password", password, "regip", regip, "regdate", time.Now().Unix(), "maxfriends", 1000, "headurl", "", "desc", "", "app", appname, "parent", puid)
	conn.Send("SADD", "fgroup:"+String(uid), defaultGroupName)
	conn.Send("SADD", "user:"+appname, uid)
	conn.Send("SADD", "app:user", uid)

	_, err = conn.Do("EXEC")

	if err != nil {
		fmt.Println("createAppUser error:", err.Error())
		return ERR_REDIS, 0
	}

	return ERR_NONE, uid
}

func (this *redisDataManager) isAppUser(appname string, uid uint64) bool {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("SISMEMBER", "user:"+appname, uid)

	if err != nil {
		fmt.Println("isAppUser error:", err.Error())
		return false
	}

	return Bool(ret)
}

func (this *redisDataManager) checkAppLogin(appname, password string) int {
	if password == "" {
		return ERR_PASSWORD_INVALID
	}
	conn := this.redisPool.Get()
	defer conn.Close()
	fmt.Println("checkAppLogin:", password, appname)
	// ret, err := conn.Do("HGET", uid, "password")

	// if err != nil {
	// 	fmt.Println("checkAppLogin error:", err.Error())
	// 	return ERR_REDIS
	// }

	// if ret == nil {
	// 	return ERR_USER_NOT_EXIST
	// }

	// if String(ret) != password {
	// 	return ERR_PASSWORD_INVALID
	// }

	// ret, err = conn.Do("SISMEMBER", "app", appname)

	// if err != nil {
	// 	fmt.Println("checkAppLogin error:", err.Error())
	// 	return ERR_REDIS
	// }

	// if ret == nil {
	// 	return ERR_APP_NOT_EXIST
	// }

	return ERR_NONE
}

func (this *redisDataManager) setAppVerifyData(uuid string, uid uint64) int {
	conn := this.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", uuid, uid, "EX", 30)

	if err != nil {
		fmt.Println("setAppVerifyData error:", err.Error())
		return ERR_REDIS
	}

	return ERR_NONE
}

func (this *redisDataManager) verifyAppLoginData(uuid string, uid uint64) bool {
	conn := this.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("GET", uuid)

	if err != nil {
		fmt.Println("setAppVerifyData error:", err.Error())
		return false
	}

	if ret == nil {
		return false
	}

	return Uint64(ret) == uid
}
