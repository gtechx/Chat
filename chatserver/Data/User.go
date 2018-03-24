package data

import (
	"time"

	"github.com/nature19862001/Chat/chatserver/Config"
	. "github.com/nature19862001/base/common"
)

func (this *RedisDataManager) CreateAccount(account, password, regip string) error {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("INCR", "UID")

	if err != nil {
		return err
	}

	uid := Uint64(ret)

	conn.Send("MULTI")
	conn.Send("HMSET", uid, "account", account, "password", password, "regip", regip, "regdate", time.Now().Unix(), "maxfriends", 1000, "headurl", "", "desc", "")
	conn.Send("SADD", "group:"+String(uid), defaultGroupName)
	conn.Send("HSET", "account:uid", account, uid)
	conn.Send("SADD", "user", uid)

	_, err = conn.Do("EXEC")
	return err
}

func (this *RedisDataManager) SetMaxFriends(uid uint64, count int) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", uid, "maxfriends", count)
	return err
}

func (this *RedisDataManager) SetDesc(uid uint64, desc string) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", uid, "desc", desc)
	return err
}

func (this *RedisDataManager) IsAccountExists(account string) (bool, error) {
	conn := this.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HEXISTS", "account:uid", account)
	return Bool(ret), err
}

func (this *RedisDataManager) IsUIDExists(uid uint64) (bool, error) {
	conn := this.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("EXISTS", uid)
	return Bool(ret), err
}

func (this *RedisDataManager) GetUIDByAccount(account string) (uint64, error) {
	conn := this.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGET", "account:uid", account)
	return Uint64(ret), err
}

func (this *RedisDataManager) GetAccountByUID(uid uint64) (string, error) {
	conn := this.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGET", uid, "account")
	return String(ret), err
}

func (this *RedisDataManager) GetPassword(uid uint64) (string, error) {
	conn := this.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGET", uid, "password")
	return String(ret), err
}

func (this *RedisDataManager) SetUserOnline(uid uint64) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	conn.Send("MULTI")
	conn.Send("HSET", uid, "online", config.ServerAddr)
	conn.Send("SADD", "online", uid)
	_, err := conn.Do("EXEC")
	return err
}

func (this *RedisDataManager) IsUserOnline(uid uint64) (bool, error) {
	conn := this.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HEXISTS", uid, "online")
	return Bool(ret), err
}

func (this *RedisDataManager) SetUserState(uid uint64, state uint8) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", uid, "state", state)
	return err
}

func (this *RedisDataManager) AddUserToBlack(uid, otheruid uint64) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SADD", "black:"+String(uid), otheruid)
	return err
}

func (this *RedisDataManager) RemoveUserFromBlack(uid, otheruid uint64) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SREM", "black:"+String(uid), otheruid)
	return err
}

func (this *RedisDataManager) IsUserInBlack(uid, otheruid uint64) (bool, error) {
	conn := this.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("SISMEMBER", "black:"+String(uid), otheruid)
	return Bool(ret), err
}
