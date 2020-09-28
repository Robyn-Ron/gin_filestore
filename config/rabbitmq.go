package config

const (
	AsyncTransferEnable = true //是否开启文件异步转移(默认同步)
	RabbitURL = "amqp://guest:guest@127.0.0.1:5672/" //rabbitmq服务的入口url
	TransExchangeName = "uploadserver.trans" //用于文件transfer的交换机
	TransOSSQueueName = "uploadserver.trans.oss" //oss转移队列名
	TransOSSErrQueueName = "uploadserver.trans.oss.err" //oss转移失败后写入另一个队列的队列名
	TransOSSRoutingKey = "oss" //routingkey
)
