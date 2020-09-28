package mq

import "log"

var done chan bool

//监听队列，获取消息
func StartConsume(qName, cName string, callback func(msg []byte) bool)  {
	//1.调用rabbitmq方法获取消息信道
	msgs, err := channel.Consume(qName, cName, true, false, false, false, nil)
	if err != nil{
		log.Println(err.Error())
		return
	}

	//2.循环从信道里面获取消息

	done = make(chan bool)
	go func() {
		for msg := range msgs {
			//3.调用callback方法来处理新的消息:msg
			if !callback(msg.Body) {
				//TODO:将任务写到另一个队列，用于异常情况的重试

			}
		}

		done <- true
	}()

	//阻塞
	<- done
}
