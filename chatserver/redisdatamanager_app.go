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
