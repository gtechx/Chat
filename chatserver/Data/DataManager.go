package data

import (
	"time"

	"github.com/garyburd/redigo/redis"
	. "github.com/nature19862001/base/common"
)

type RedisDataManager struct {
	redisPool *redis.Pool
}

var instanceDataManager *RedisDataManager

func Manager() *RedisDataManager {
	if instanceDataManager == nil {
		instanceDataManager = &RedisDataManager{}
	}
	return instanceDataManager
}

func (this *RedisDataManager) Initialize() error {
	this.redisPool = &redis.Pool{
		MaxIdle:      3,
		IdleTimeout:  240 * time.Second,
		Dial:         redisDial,
		TestOnBorrow: redisOnBorrow,
	}

	conn := this.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("EXISTS", "UID")

	if err != nil {
		return err
	}

	if !Bool(ret) {
		_, err = conn.Do("SET", "UID", 10000)

		if err != nil {
			return err
		}
	}

	ret, err = conn.Do("HEXISTS", "admin", 0)

	if err != nil {
		return err
	}

	if !Bool(ret) {
		_, err = conn.Do("HSET", "admin", 0, 0xffffffff)

		if err != nil {
			return err
		}

		_, err = conn.Do("HSET", 0, "password", Md5("ztgame@123"))

		if err != nil {
			return err
		}
	}

	return err
}

func redisDial() (redis.Conn, error) {
	c, err := redis.Dial("tcp", redisAddr)
	if err != nil {
		return nil, err
	}
	if _, err := c.Do("AUTH", "ztgame@123"); err != nil {
		c.Close()
		return nil, err
	}
	// if _, err := c.Do("SELECT", db); err != nil {
	// 	c.Close()
	// 	return nil, err
	// }
	return c, nil
}

func redisOnBorrow(c redis.Conn, t time.Time) error {
	if time.Since(t) < time.Minute {
		return nil
	}
	_, err := c.Do("PING")
	return err
}
