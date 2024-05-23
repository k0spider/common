#### redis基础组件
不支持redis7及以上
##### 配置
> ```
> redis:
>   host: "127.0.0.1"
>   port: "6379"
>   auth: "pwd"
>   userName: "username"
>   db: 0
>   poolSize: 10
>   encryption: 0
>   framework: ""  # 如果为集群该值为"cluster"
> ```
##### 读取配置
> ```
> import "github.com/k0spider/common/dao/redis"
> 
> type Config struct {
>   Redis   *redis.RedisConfig     `yaml:"redis"`
> }
> ```
###### 初始化
> ```
> redis.InitRedis(Config.Redis)
> ```

###### 使用
> ```
> val, err := Redis.Get(context.Background(), "xxx").Result()
> ```