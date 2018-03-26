package cdata

import (
	"time"

	"github.com/nature19862001/Chat/chatserver/Config"
	"github.com/nature19862001/Chat/chatserver/Entity"
	. "github.com/nature19862001/base/common"
)

//每个app之间可以是独立的数据，也可以共享数据，根据你的设置
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
	conn.Send("HSET", "account:uid", account, uid)
	conn.Send("SADD", "user", uid)

	_, err = conn.Do("EXEC")
	return err
}

func (this *RedisDataManager) CreateAppData(entity centity.UserEntity) error {
	conn := this.redisPool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("HSET", entity.KeyAppData, "createdate", time.Now().Unix())
	conn.Send("SADD", entity.KeyGroup, defaultGroupName)
	_, err = conn.Do("EXEC")
	return err
}

func (this *RedisDataManager) DeleteAppData(entity centity.UserEntity) error {
	conn := this.redisPool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("HDEL", entity.KeyAppData)
	conn.Send("DEL", entity.KeyGroup)
	conn.Send("HDEL", entity.KeyFriend)
	conn.Send("HDEL", entity.KeyFriendRequest)
	conn.Send("HDEL", entity.KeyBlack)
	_, err = conn.Do("EXEC")
	return err
}

func (this *RedisDataManager) IsAppDataExists(entity centity.UserEntity) (bool, error) {
	conn := this.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HEXISTS", entity.KeyAppData)
	return Bool(ret), err
}

func (this *RedisDataManager) SetAppDataConfig(entity centity.UserEntity, configname string, data interface{}) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", entity.KeyAppData, configname)
	return err
}

func (this *RedisDataManager) GetAppDataConfig(entity centity.UserEntity, configname string) (interface{}, error) {
	conn := this.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGET", entity.KeyAppData, configname)
	return ret, err
}

// func (this *RedisDataManager) SetMaxFriends(uid uint64, count int) error {
// 	conn := this.redisPool.Get()
// 	defer conn.Close()
// 	_, err := conn.Do("HSET", uid, "maxfriends", count)
// 	return err
// }

// func (this *RedisDataManager) SetDesc(uid uint64, desc string) error {
// 	conn := this.redisPool.Get()
// 	defer conn.Close()
// 	_, err := conn.Do("HSET", uid, "desc", desc)
// 	return err
// }

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

func (this *RedisDataManager) SetUserOnline(entity centity.UserEntity) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	conn.Send("MULTI")
	conn.Send("HSET", entity.KeyAppData, "online", config.ServerAddr)
	conn.Send("SADD", "online", entity.UID())
	_, err := conn.Do("EXEC")
	return err
}

func (this *RedisDataManager) IsUserOnline(entity centity.UserEntity) (bool, error) {
	conn := this.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HEXISTS", entity.KeyAppData, "online")
	return Bool(ret), err
}

func (this *RedisDataManager) SetUserState(entity centity.UserEntity, state uint8) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", entity.KeyAppData, "state", state)
	return err
}

func (this *RedisDataManager) AddUserToBlack(entity centity.UserEntity, otheruid uint64) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SADD", entity.KeyBlack, otheruid)
	return err
}

func (this *RedisDataManager) RemoveUserFromBlack(entity centity.UserEntity, otheruid uint64) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SREM", entity.KeyBlack, otheruid)
	return err
}

func (this *RedisDataManager) IsUserInBlack(entity centity.UserEntity, otheruid uint64) (bool, error) {
	conn := this.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("SISMEMBER", entity.KeyBlack, otheruid)
	return Bool(ret), err
}
