package data

import (
	//"errors"
	"fmt"

	"github.com/garyburd/redigo/redis"
	. "github.com/nature19862001/base/common"
)

func (this *RedisDataManager) pullMsg(addr string, timeout int) []byte {
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
		retarr, _ := redis.Values(ret, nil)
		//fmt.Println(err.Error())
		return Bytes(retarr[1])
	}
}

func (this *RedisDataManager) sendMsgToUser(uid uint64, data []byte) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	//check if friend is exist
	ret, err := conn.Do("EXISTS", uid)

	if err != nil {
		fmt.Println("sendMsgToUser error:", err.Error())
		return ERR_REDIS
	}

	if !Bool(ret) {
		return ERR_USER_NOT_EXIST
	}

	ret, err = conn.Do("HGET", uid, "online")

	if err != nil {
		fmt.Println("sendMsgToUser error:", err.Error())
		return ERR_REDIS
	}

	data = append(Bytes(uid), data...)
	//fmt.Println(data)
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

func (this *RedisDataManager) sendMsgToRoom() {

}
