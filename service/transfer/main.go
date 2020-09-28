package main

import (
	"CloudWebOfGin/config"
	"CloudWebOfGin/mq"
	"encoding/json"
	"log"
	"os"
)

func main()  {
	mq.StartConsume(
		config.TransOSSQueueName,
		"transfer_oss",
		ProcessTransfer,
	)
}

//处理文件转移的真正逻辑
func ProcessTransfer(msg []byte) bool {
	//1.解析msg

	pubData := mq.TransferData{}
	err := json.Unmarshal(msg, pubData)
	if err != nil{
		log.Println(err.Error())
		return false
	}

	//2.根据临时存储文件路径，创建文件句柄
	filed, err := os.Open(pubData.CurLocation)
	if err != nil{
		log.Println(err.Error())
		return false
	}

	//3.通过文件句柄将文件内容读出来并且上传到OSS
	oss.

	//4.更新文件的存储路径到文件表

	return true
}