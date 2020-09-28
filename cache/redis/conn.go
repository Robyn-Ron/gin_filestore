package redis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

//定义全局连接池对象
var (
	pool *redis.Pool
	redisHost = "192.168.243.188" //IP地址
	redisPass = "" //redis密码
)

//newRedisPool:创建redis连接池
func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle: 50, //最大连接数
		MaxActive: 30, //最大活跃连接数
		IdleTimeout: 300 * time.Second, //超时时间（五分钟）
		Dial: func() (redis.Conn, error) {
			//1.打开连接
			c, err := redis.Dial("tcp", redisHost)
			if err != nil{
				fmt.Println(err)
				return nil, err
			}

			//2.访问认证.（如果redis添加了登录密码的验证方式）

			return c, nil
		},

		//健康检查
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute{
				return nil
			}

			_, err := c.Do("PING")
			if err != nil{
				return err
			}

			return nil
		},
	}
}

func init()  {
	pool = newRedisPool()
}

func RedisPool() *redis.Pool {
	return pool
}