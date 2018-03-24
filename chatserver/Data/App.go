package data

import (
	"time"
	
	. "github.com/nature19862001/base/common"
)

//[set]app aid set
//[hset]app:aid:uid aid owner desc regdate
//[hset]app:aid:uid:config

//app op
func (this *RedisDataManager) CreateApp(uid uint64, name string) error {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("INCR", "UID")

	if err != nil {
		return err
	}

	appid := Uint64(ret)

	conn.Send("MULTI")
	conn.Send("SADD", "app", appid)
	conn.Send("SADD", "app:"+String(uid), appid)
	conn.Send("HMSET", "app:"+String(uid)+":"+String(appid), "owner", uid, "desc", "", "iconurl", "", "regdate", time.Now().Unix())

	_, err = conn.Do("EXEC")

	return err
}

func (this *RedisDataManager) DeleteApp(uid, appid uint64) error {
	conn := this.redisPool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("SREM", "app", appid)
	conn.Send("SREM", "app:"+String(uid), appid)
	conn.Send("DEL", "app:"+String(uid)+":"+String(appid))

	_, err := conn.Do("EXEC")

	return err
}

func (this *RedisDataManager) IsAppExists(appid uint64) (bool, error) {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("SISMEMBER", "app", appid)

	return Bool(ret), err
}


