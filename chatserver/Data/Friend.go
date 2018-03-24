package data

import (
	"github.com/garyburd/redigo/redis"
	. "github.com/nature19862001/base/common"
)

var defaultGroupName string = "我的好友"
var userOnlineKeyName string = "user:online"

func (this *RedisDataManager) AddFriendRequest(uid, otheruid uint64, group string) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", "friend:request:"+String(uid), otheruid, group)
	return err
}

func (this *RedisDataManager) RemoveFriendRequest(uid, otheruid uint64) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HDEL", "friend:request:"+String(uid), otheruid)
	return err
}

func (this *RedisDataManager) AddFriend(uid, otheruid uint64, group string) error {
	conn := this.redisPool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("HSET", "friend:"+String(uid), otheruid, group)
	conn.Send("SADD", "group:"+String(uid)+":"+group, otheruid)
	_, err := conn.Do("EXEC")

	return err
}

func (this *RedisDataManager) RemoveFriend(uid, otheruid uint64, group string) error {
	conn := this.redisPool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("HDEL", "friend:"+String(uid), otheruid)
	conn.Send("SREM", "group:"+String(uid)+":"+group, otheruid)
	_, err := conn.Do("EXEC")

	return err
}

func (this *RedisDataManager) GetFriendList(uid uint64, group string) ([]uint64, error) {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("SMEMBERS", "group:"+String(uid)+":"+group)

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

func (this *RedisDataManager) IsFriend(uid, otheruid uint64) (bool, error) {
	conn := this.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HEXISTS", "friend:"+String(uid), otheruid)
	return Bool(ret), err
}

func (this *RedisDataManager) GetGroupOfFriend(uid, otheruid uint64) (string, error) {
	conn := this.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGET", "friend:"+String(uid), otheruid)
	return String(ret), err
}

func (this *RedisDataManager) AddGroup(uid uint64, group string) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SADD", "group:"+String(uid), group)
	return err
}

func (this *RedisDataManager) RemoveGroup(uid uint64, group string) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SREM", "group:"+String(uid), group)
	return err
}

func (this *RedisDataManager) IsGroupExists(uid uint64, group string) (bool, error) {
	conn := this.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("SISMEMBER", "group:"+String(uid), group)
	return Bool(ret), err
}

func (this *RedisDataManager) IsFriendInGroup(uid, otheruid uint64, group string) (bool, error) {
	conn := this.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("SISMEMBER", "group:"+String(uid)+":"+group, otheruid)
	return Bool(ret), err
}

func (this *RedisDataManager) MoveFriendToGroup(uid, otheruid uint64, srcgroup, destgroup string) error {
	conn := this.redisPool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("SREM", "group:"+String(uid)+":"+srcgroup, otheruid)
	conn.Send("SADD", "group:"+String(uid)+":"+destgroup, otheruid)
	_, err := conn.Do("EXEC")

	return err
}

// func (this *RedisDataManager) BanFriend(uid, fuid uint64) {

// }

// func (this *RedisDataManager) UnBanFriend(uid, fuid uint64) {

// }

func (this *RedisDataManager) SetFriendVerifyType(uid uint64, vtype byte) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", uid, "verifytype", vtype)
	return err
}

func (this *RedisDataManager) GetFriendVerifyType(uid uint64) (byte, error) {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("HGET", uid, "verifytype")

	return Byte(ret), err
}
