package cdata

import (
	"github.com/garyburd/redigo/redis"
	"github.com/nature19862001/Chat/chatserver/Entity"
	. "github.com/nature19862001/base/common"
)

var defaultGroupName string = "我的好友"
var userOnlineKeyName string = "user:online"

func (this *RedisDataManager) AddFriendRequest(entity centity.UserEntity, otheruid uint64, group string) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", entity.KeyFriendRequest, otheruid, group)
	return err
}

func (this *RedisDataManager) RemoveFriendRequest(entity centity.UserEntity, otheruid uint64) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HDEL", entity.KeyFriendRequest, otheruid)
	return err
}

func (this *RedisDataManager) AddFriend(entity centity.UserEntity, otheruid uint64, group string) error {
	conn := this.redisPool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("HSET", entity.KeyFriend, otheruid, group)
	conn.Send("SADD", entity.KeyGroup+":"+group, otheruid)
	_, err := conn.Do("EXEC")

	return err
}

func (this *RedisDataManager) RemoveFriend(entity centity.UserEntity, otheruid uint64, group string) error {
	conn := this.redisPool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("HDEL", entity.KeyFriend, otheruid)
	conn.Send("SREM", entity.KeyGroup+":"+group, otheruid)
	_, err := conn.Do("EXEC")

	return err
}

func (this *RedisDataManager) GetFriendList(entity centity.UserEntity, group string) ([]uint64, error) {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("SMEMBERS", entity.KeyGroup+":"+group)

	if err != nil {
		return nil, err
	}

	retarr, err := redis.Values(ret, nil)

	if err != nil {
		return nil, err
	}

	userlist := []uint64{}
	for _, uid := range retarr {
		userlist = append(userlist, Uint64(uid))
	}

	return userlist, err
}

func (this *RedisDataManager) IsFriend(entity centity.UserEntity, otheruid uint64) (bool, error) {
	conn := this.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HEXISTS", entity.KeyFriend, otheruid)
	return Bool(ret), err
}

func (this *RedisDataManager) GetGroupOfFriend(entity centity.UserEntity, otheruid uint64) (string, error) {
	conn := this.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGET", entity.KeyFriend, otheruid)
	return String(ret), err
}

func (this *RedisDataManager) AddGroup(entity centity.UserEntity, group string) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SADD", entity.KeyGroup, group)
	return err
}

func (this *RedisDataManager) RemoveGroup(entity centity.UserEntity, group string) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SREM", entity.KeyGroup, group)
	return err
}

func (this *RedisDataManager) IsGroupExists(entity centity.UserEntity, group string) (bool, error) {
	conn := this.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("SISMEMBER", entity.KeyGroup, group)
	return Bool(ret), err
}

func (this *RedisDataManager) IsFriendInGroup(entity centity.UserEntity, otheruid uint64, group string) (bool, error) {
	conn := this.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("SISMEMBER", entity.KeyGroup+":"+group, otheruid)
	return Bool(ret), err
}

func (this *RedisDataManager) MoveFriendToGroup(entity centity.UserEntity, otheruid uint64, srcgroup, destgroup string) error {
	conn := this.redisPool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("SREM", entity.KeyGroup+":"+srcgroup, otheruid)
	conn.Send("SADD", entity.KeyGroup+":"+destgroup, otheruid)
	_, err := conn.Do("EXEC")

	return err
}

// func (this *RedisDataManager) BanFriend(uid, fuid uint64) {

// }

// func (this *RedisDataManager) UnBanFriend(uid, fuid uint64) {

// }

// func (this *RedisDataManager) SetFriendVerifyType(uid uint64, vtype byte) error {
// 	conn := this.redisPool.Get()
// 	defer conn.Close()
// 	_, err := conn.Do("HSET", uid, "verifytype", vtype)
// 	return err
// }

// func (this *RedisDataManager) GetFriendVerifyType(uid uint64) (byte, error) {
// 	conn := this.redisPool.Get()
// 	defer conn.Close()

// 	ret, err := conn.Do("HGET", uid, "verifytype")

// 	return Byte(ret), err
// }
