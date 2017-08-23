package main

import (
	//"errors"
	"fmt"
	//"github.com/garyburd/redigo/redis"
	. "github.com/nature19862001/base/common"
	//"strings"
	//"time"
)

func (this *redisDataManager) checkAppLogin(uid uint64, password, appname string) int {
	conn := this.redisPool.Get()
	defer conn.Close()
	fmt.Println("checkAppLogin:", uid, password, appname)
	ret, err := conn.Do("HGET", uid, "password")

	if err != nil {
		fmt.Println("checkAppLogin error:", err.Error())
		return ERR_REDIS
	}

	if ret == nil {
		return ERR_USER_NOT_EXIST
	}

	if String(ret) != password {
		return ERR_PASSWORD_INVALID
	}

	ret, err = conn.Do("SISMEMBER", "app", appname)

	if err != nil {
		fmt.Println("checkAppLogin error:", err.Error())
		return ERR_REDIS
	}

	if ret == nil {
		return ERR_APP_NOT_EXIST
	}

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
