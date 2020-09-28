package mq

import (
	"CloudWebOfGin/config"
	"github.com/streadway/amqp"
	"log"
)

var conn *amqp.Connection

var channel *amqp.Channel

func initChannel() bool {
	//1.判断channel是否已经创建过
	if channel != nil{
		return true
	}

	//2.获得rabbitmq的一个连接
	conn, err := amqp.Dial(config.RabbitURL)
	if err != nil{
		log.Println(err.Error())
		return false
	}

	//3.打开一个channel，用于消息的发布与接收
	channel, err := conn.Channel()
	if err != nil{
		log.Println(err.Error())
		return false
	}

	defer channel.Close()

	return true
}

//发布消息
func Publish(exchange, routingKey string, msg []byte) bool {
	// 1. 发布消息之前先检查channel是否是正常的
	if !initChannel(){
		return false
	}
	// 2. 执行消息发布动作
	err := channel.Publish(exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "text/plain", //明文格式
		Body: msg, //信息
	})

	if err != nil{
		log.Println(err.Error())
		return false
	}

	return true
}

