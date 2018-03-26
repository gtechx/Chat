package cdata

import (
	"time"

	"github.com/garyburd/redigo/redis"
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
	conn.Send("HMSET", "app:"+String(appid), "owner", uid, "desc", "", "iconurl", "", "regdate", time.Now().Unix())

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

func (this *RedisDataManager) AddAppZone(appid uint64, zones ...uint32) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGET", "app:zone:"+String(appid), zones)
}

func (this *RedisDataManager) GetAppOwner(appid uint64) (uint64, error) {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("HGET", "app:"+String(appid), "owner")

	return Uint64(ret), err
}

func (this *RedisDataManager) IsAppZone(appid uint64, zone uint32) (bool, error) {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("SISMEMBER", "app:zone:"+String(appid), zone)

	return Bool(ret), err
}

func (this *RedisDataManager) AddShareApp(uid, appid, otherappid uint64) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	conn.Send("MULTI")
	conn.Send("HMSET", "app:"+String(uid)+":"+String(appid), "share", otherappid)
	conn.Send("SADD", "app:share:"+String(otherappid), appid)
	_, err := conn.Do("EXEC")
	return err
}

func (this *RedisDataManager) IsShareWithOtherApp(uid, appid uint64) (bool, error) {
	conn := this.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HEXISTS", "app:"+String(uid)+":"+String(appid), "share")
	return Bool(ret), err
}

func (this *RedisDataManager) GetShareApp(appid uint64) (uint64, error) {
	conn := this.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGET", "app:"+String(uid)+":"+String(appid), "share")
	return Uint64(ret), err
}

func (this *RedisDataManager) GetMyShareAppList(appid uint64) ([]uint64, error) {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("SMEMBERS", "app:share:"+String(appid))

	if err != nil {
		return nil, err
	}

	retarr, err := redis.Values(ret, nil)

	if err != nil {
		return nil, err
	}

	applist := []uint64{}
	for _, otherappid := range retarr {
		applist = append(applist, Uint64(otherappid))
	}

	return applist, err
}
