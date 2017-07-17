package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	. "github.com/nature19862001/base/common"
)

var RedisConn redis.Conn

func InitRedis() error {
	var err error
	RedisConn, err = redis.Dial("tcp", "127.0.0.1:6379")

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
