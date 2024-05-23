#### AWS文件上传
##### 依赖
无
##### 配置
> ```
> awsUpload:
>   accessKeyID: "accessKeyID"
>   secretAccessKey: "secretAccessKey"
>   bucket: "bucket"
>   region: "region"
>   folder: "folder"
> ```
##### 读取配置
> ```
> import "github.com/k0spider/common/upload/aws"
> 
> type Config struct {
>   AwsUpload  aws.Config   `yaml:"awsUpload"`
> }
> ```

###### 获取AWS文件上传预签名URL
> ```
> // 获取到AWS文件上传预签名URL后 前端使用该URL发起PUT请求向AWS提交文件「注意：header内容在PUT请求中必须携带」
> func PresignUpload(c echo.Context) error {
> 	ctx := c.Request().Context()
> 	req := aws.AwsS3Request{}
> 	err := c.Bind(&req)
> 	if err != nil {
> 		return code.ErrLog(ctx, defined.ParamsErr, err)
> 	}
> 	ext := path.Ext(req.FileName)
> 	if ext == "" {
> 		return code.ErrLog(ctx, defined.ParamsErr, nil)
> 	}
> 	fileName := fmt.Sprintf("%s/%s%s", time.Now().Format(utils.DateFolder), utils.GenerateToken(32), ext)
> 	res, err := aws.PresignPutObject(Config.AwsUpload, fileName, 300)
> 	if err != nil {
> 		return code.ErrLog(ctx, defined.FailedErr, err)
> 	}
> 	return c.JSON(http.StatusOK, code.Success(ctx, res))
> }
> 
> ```