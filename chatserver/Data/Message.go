package cdata

import (
	//"errors"

	"github.com/garyburd/redigo/redis"
	. "github.com/nature19862001/base/common"
)

func (this *RedisDataManager) PullOnlineMessage(serveraddr string, timeout int) ([]byte, error) {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("BLPOP", "message:"+serveraddr, timeout)

	if err != nil {
		return nil, err
	}

	retarr, err := redis.Values(ret, nil)

	if err != nil {
		return nil, err
	}

	return Bytes(retarr[1]), err
}

func (this *RedisDataManager) GetOfflineMessage(uid, appid uint64) ([][]byte, error) {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("LRANGE", "message:offline:"+String(uid)+":"+String(appid), 0, -1)

	if err != nil {
		return nil, err
	}

	retarr, err := redis.Values(ret, nil)

	if err != nil {
		return nil, err
	}

	msglist := [][]byte{}
	for i := 1; i < len(retarr); i++ {
		msglist := append(msglist, Bytes(retarr[i]))
	}

	return msglist, err
}

func (this *RedisDataManager) SendMsgToUserOnline(uid, appid uint64, data []byte, serveraddr string) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("RPUSH", "message:"+serveraddr, data)
	return err
}

func (this *RedisDataManager) SendMsgToUserOffline(uid, appid uint64, data []byte) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("RPUSH", "message:offline:"+String(uid)+":"+String(appid), data)
	return err
}

func (this *RedisDataManager) SendMsgToRoom() {

}
