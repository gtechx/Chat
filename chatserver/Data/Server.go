package data

import (
	//"errors"

	"github.com/garyburd/redigo/redis"
	. "github.com/nature19862001/base/common"
)

//[sorted sets]serverlist pair(count,addr)
//[sets]ttl:addr
var serverListKeyName string = "serverlist"

//server op
func (this *RedisDataManager) RegisterServer(addr string) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("ZADD", serverListKeyName, 0, addr)

	return err
}

func (this *RedisDataManager) IncrByServerClientCount(addr string, count int) error {
	conn := this.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("ZINCRBY", serverListKeyName, count, addr)

	return err
}

func (this *RedisDataManager) GetServerList() ([]string, error) {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("ZRANGE", serverListKeyName, 0, -1)

	if err != nil {
		return nil, err
	}

	slist, _ := redis.Strings(ret, err)
	return slist, err
}

func (this *RedisDataManager) GetServerCount() (int, error) {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("ZCARD", serverListKeyName)

	return Int(ret), err
}

func (this *RedisDataManager) SetServerTTL(addr string, seconds int) error {
	conn := this.redisPool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", "ttl:"+addr, 0, "EX", seconds)

	return err
}

func (this *RedisDataManager) CheckServerTTL() error {
	return nil
}

func (this *RedisDataManager) VoteServerDie() error {
	return nil
}
