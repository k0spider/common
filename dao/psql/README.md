### psql基础组件

#### Gorm

##### 配置

> ```
> gormCfg:
>   driver: postgres
>   debug: true
>   dsn:                    # gorm 不支持多节点连接，这里默认连接第一个dsn
>       - "postgres://postgres:pwd@127.0.0.1:55001/admin?sslmode=disable"
>   maxOpens: 20            # 设置池上的最大连接数
>   maxIdles: 10            # 设置池上的最大空闲连接数
>   idleTimeout: 30         # 设置可以重复使用连接的最大时间s
>   logDir: "./logs/gorm"  # 设置sql日志报错位置
>   logMaxDay: 30  # 设置sql日志最大保存时长 天
> ```

##### 读取配置

> ```
> import "github.com/k0spider/common/dao/psql"
> 
> type Config struct {
>   GormCfg     *psql.GormConfig   `yaml:"gormCfg"`
> }
> ```

###### 初始化

> ```
> GormDb := psql.NewGorm(Config.GormCfg)
> ```

###### 事务使用

> ```
> 使用gorm库的事务操作
> ```