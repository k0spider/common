#### MQ基础组件
##### 依赖
1.  依赖日志组件「需要日志配置」
##### 配置
> ```
> rabbitMQ:
>   dsn: "amqp://root:root@127.0.0.1:5672//admin"
>   debug: true
>   notify:
>       exchangeName: "notify-exchange"
>       queueName: "notify"
>       routeKey: "notify-key"
>   subscribe:
>       exchangeName: "subscribe-exchange"
>       queueName: "subscribe"
>       routeKey: "subscribe-key"
>   timeOut:
>       exchangeName: "time-out-exchange"
>       queueName: "time-out"
>       routeKey: "time-out-key"
> ```
##### 读取配置
> ```
> import "github.com/k0spider/common/dao/mq"
> 
> type Config struct {
>   RabbitMQ    *mq.RabbitMQPoolConfig `yaml:"rabbitMQ"`
> }
> ```
###### 初始化
> ```
> 
> mq.InitPool(Config.RabbitMQ)
> ```

###### 延迟队列使用
> ```
> err = mq.Push(ctx, Config.RabbitMQ.Notify, &amqp.Publishing{
> 		Body:       []byte("test"),
> 		Expiration: strconv.Itoa(int(1 * 1000)),
> })
> 
> ```
###### 常规队列使用
> ```
>  err = mq.Push(ctx, Config.RabbitMQ.Subscribe, &amqp.Publishing{
> 		Body:  []byte("test"),
>  })
> ```
###### 监听消息
> ```
> // 允许向Consumers传递lark配置在发生异常是发送警报
>  err := mq.Consumers(ctx, Config.RabbitMQ.TimeOut, func(delivery amqp.Delivery) {
>		// 监听处理逻辑
>		delivery.Ack(false)
>  })
> ```

###### 获取监听消息channel
> ```
>  channelMsg, err := mq.GetMsgChannel(Config.RabbitMQ.TimeOut)
>  if err != nil {
>   panic(err)
>  }
>  for delivery := range channelMsg {
>       // 处理逻辑
>       delivery.Ack(false)
> }
> ```