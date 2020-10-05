package handler

import (
	myredis "file_store_net_http/cache/redis"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"runtime"
	"testing"
)

func TestRedisMulti(t *testing.T) {
	conn := myredis.RedisPool().Get()
	defer conn.Close()

	conn.Send("MULTI")
	//conn.Send("HSET", "redis_test", "fieldOne", 11)
	//conn.Send("HSET", "redis_test", "fieldTwo", 22)
	//conn.Send("HSET", "redis_test", "fieldThree", 33)
	conn.Send("HGETALL", "redis_test")

//这里的data是一个[]interface{}类型, 其中每一个下标对应于send语句的出现位置;
//		而对于只有一个的redis.DO()操作, 则直接使用data[0].([]interface{}/int64)进行类型断言;

//而每一个Send的返回结果, 如果是查询类, 断言为data[0].([]interface{})类型, 而且[]interface{}中的
//		每一个成员都是[]byte类型;
//	如果是set等类, 断言为data[0].(int64)类型;
	//这里的values的返回值就是一个interface{}类型;
	data, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		_, fn, line, _ := runtime.Caller(0)
		log.Println("fliename:", fn, "; line:", line, "; error:",err)
		return
	}
	fmt.Println("_________________")
	val := data[0].([]interface{})
	for _, unit := range val{
		fmt.Println(string(unit.([]byte)))
	}
	//unit := value.(string)
	//fmt.Printf("index=%v, value%v, unit=%v\n", index, value, unit)
}

