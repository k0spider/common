#### 邮件发送基础组件
##### 组件能力
1. 支持Gmail邮件发送
2. 支持QQ邮件发送
##### 配置
> ```
> email:
>   identity: ""
>   userName: "user name" # 账号
>   password: "password" # 授权码「通常为16位小写字母组成」
> ```
##### 读取配置
> ```
> import "github.com/k0spider/common/email"
> 
> type Config struct {
>   Email    *email.Config `yaml:"email"`
> }
> ```

###### 参数
> ```
> type EmailParameter struct {
> 	To                 []string // 收件人
> 	Bcc                []string // 抄送人
> 	Cc                 []string // 抄送人
> 	Subject            string   // 自定义邮件标题类容「如果定义优先使用」
> 	TextBody           []byte   // 邮件文本类容「如果定义优先使用」
> 	HtmlBody           []byte   // HTML邮件内容
> }
> ```

###### 发送
> ```
> err := NewEmail(Cnfig.Email).Gmail().Send(&email.EmailParameter{
>       To: []string{"xxx@qq.com.com"},
>       Subject: "test",
>       TextBody: []byte("code:121212"),
> })
> fmt.Println(err)
> ```