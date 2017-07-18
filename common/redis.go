package common

import (
	"github.com/garyburd/redigo/redis"
	//. "github.com/nature19862001/base/common"
)

var RedisConn redis.Conn

func InitRedis(net, addr string) error {
	var err error
	RedisConn, err = redis.Dial(net, addr)

	if err != nil {
		return err
	}

	return nil
}

func CloseRedis() {
	if RedisConn != nil {
		RedisConn.Close()
		RedisConn = nil
	}
}
