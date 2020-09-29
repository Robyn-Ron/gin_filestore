package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

//返回前端响应体数据类型: 一般需要前后端讨论, 制定统一接口文档, 再自定义返回结构体map, 或组装返回结构体map;

type RespMes struct{
	Code int 			`json:"code"`
	Msg string 			`json:"msg"`
	Data interface{} 	`json:"data"`
}

func NewRespMsg(code int, msg string, data interface{}) *RespMes{
	return &RespMes{
		Code: code,
		Msg: msg,
		Data: data,
	}
}

func (this *RespMes)JsonBytes() []byte{
	data, err := json.Marshal(this)
	if err != nil {
		log.Println(err)
	}
	return data
}

func (this *RespMes)JsonString() string {
	data, err := json.Marshal(this)
	if err != nil {
		log.Println(err)
	}
	return string(data)
}

func GenSimpleRespStream(code int, msg string)[] byte{
	return []byte(fmt.Sprintf("{code:%s,msg:%s}", strconv.Itoa(code), msg))
}

func GenSimpleRespString(code int, msg string) string{
	return fmt.Sprintf("{code:%s,msg:%s}", strconv.Itoa(code), msg)
}