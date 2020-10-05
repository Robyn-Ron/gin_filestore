package myredis

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"runtime"
	"time"
)

var (
	pool *redis.Pool
)

func initRedisPool() *redis.Pool {
	//Pool只是一个抽象出来的pool, 其中管理着的还是一个个redis.Conn的连接;
	p := &redis.Pool{
		Dial:	func()(redis.Conn, error){
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				_, fn, line, _ := runtime.Caller(0)
				log.Println("filename: ", fn, "; line:", line, "; err:", err)
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow:    func(c redis.Conn, t time.Time)error{ //定时任务, 检查conn与redis-server的连接
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			if err != nil {
				_, fn, line, _ := runtime.Caller(0)
				log.Println("filename: ", fn, "; line:", line, "; err:", err)
				return err
			}
			return nil
		},
		MaxIdle:         50,
		MaxActive:       30,
		IdleTimeout:     300 * time.Second,
	}
	return p
}

func init() {
	pool = initRedisPool()
}

func RedisPool() *redis.Pool {
	return pool
}