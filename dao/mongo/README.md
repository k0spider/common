#### mongo驱动集成

##### 配置
> ```
> projectName:
>   appName: "xxx-projectName" # 您的项目名称
>   # 遵循URI格式 https://www.mongodb.com/docs/manual/reference/connection-string/
>   dns: "mongodb://username:password@localhost:27017"
>   database: "test" # 库名 
>   maxOpens: 20    # 设置池上的最大连接数
>   minOpens: 10    # 设置池上的最小维持连接数
>   idleTimeout: 30 # 设置可以重复使用连接的最大时间s
> ```
##### 读取配置
> ```
> import "github.com/k0spider/common/dao/mongo"
> 
> type Config struct {
>   MongoCfg    *mongo.MongoConfig   `yaml:"mongoCfg"`
> }
> ```
###### 初始化
> ```
> MongoDb := mongo.NewMongo(Config.MongoCfg)
> ```

###### 使用
> ```
> fmt.Println(MongoDb.Collection("test").Name())
> ```

###### 连接池版本
```go

testClient := NewMongoClient(c)
// 从连接池创建新的会话
session, db := testClient.GetConn()
// 关闭会话，会话回到连接池
defer testClient.CloseConn(session)

```