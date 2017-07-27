package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	. "github.com/nature19862001/base/common"
)

var serverListKeyName string = "serverlist"

type redisDataManager struct {
	var redisPool *redis.Pool
}

func (this *redisDataManager) initialize() {
	this.redisPool = &redis.Pool{
		MaxIdle:      3,
		IdleTimeout:  240 * time.Second,
		Dial:         redisDial,
		TestOnBorrow: redisOnBorrow,
	}
}

func redisDial() (redis.Conn, error) {
	c, err := redis.Dial("tcp", server)
	if err != nil {
		return nil, err
	}
	if _, err := c.Do("AUTH", password); err != nil {
		c.Close()
		return nil, err
	}
	if _, err := c.Do("SELECT", db); err != nil {
		c.Close()
		return nil, err
	}
	return c, nil
}

func redisOnBorrow(c redis.Conn, t time.Time) error {
	if time.Since(t) < time.Minute {
		return nil
	}
	_, err := c.Do("PING")
	return err
}

//server op
func (this *redisDataManager) registerServer(addr string) bool {
	conn := this.redisPool.Get()
	defer conn.Close()
	n, err := conn.Do("ZADD", serverListKeyName, 0, addr)

	if err != nil {
		fmt.Println("register server error:", err.Error())
		return false
	}

	return true
}

func (this *redisDataManager) incrServerClientCountBy(addr string, count int) {
	conn := this.redisPool.Get()
	defer conn.Close()
	n, err := conn.Do("ZINCRBY", serverListKeyName, count, addr)

	if err != nil {
		fmt.Println("incrServerClientCountBy error:", err.Error())
	}
}

func (this *redisDataManager) getServerList() []string {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("ZRANGE", serverListKeyName, 0, -1)

	if err != nil {
		fmt.Println("getServerList error:", err.Error())
		return []string{}
	}

	slist, _ = redis.Strings(ret, err)
	return slist
}

func (this *redisDataManager) getServerCount() int{
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("ZCARD", serverListKeyName)

	if err != nil {
		fmt.Println("getServerCount error:", err.Error())
		return 0
	}

	return Int(ret)
}

func (this *redisDataManager) setServerTTL(addr string, seconds int) {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("SET", "ttl"+addr, 0, "EX", seconds)

	if err != nil {
		fmt.Println("setServerTTL error:", err.Error())
		return
	}
}

func (this *redisDataManager) checkServerTTL() int {
	
}

func (this *redisDataManager) voteServerDie() int {
	
}

func (this *redisDataManager) pullMsg(addr string, timeout int) []byte {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("BLPOP", "msg"+addr, timeout)

	if err != nil {
		fmt.Println("pullMsg error:", err.Error())
		return nil
	}

	if ret == nil {
		return nil
	}else{
		return Bytes(ret)
	}
}