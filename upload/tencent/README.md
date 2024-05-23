#### 腾讯Cos文件上传
##### 依赖
无
##### 配置
> ```
> cosUpload:
>   accessKeyID: "accessKeyID"
>   secretAccessKey: "secretAccessKey"
>   bucket: "bucket"
>   region: "region"
>   folder: "folder"
>   showUrl: "showUrl"
> ```
##### 读取配置
> ```
> import "github.com/k0spider/common/upload/tencent"
> 
> type Config struct {
>   CosUpload  cos.Config   `yaml:"cosUpload"`
> }
> ```

###### 获取COS文件上传预签名URL
> ```
> // 获取到TENCENT文件上传预签名URL后 前端使用该URL发起PUT请求向TENCENT提交文件
> func PresignUpload(c echo.Context) error {
> 	ctx := c.Request().Context()
> 	req := tencent.CosRequest{}
> 	err := c.Bind(&req)
> 	if err != nil {
> 		return code.ErrLog(ctx, defined.ParamsErr, err)
> 	}
> 	ext := path.Ext(req.FileName)
> 	if ext == "" {
> 		return code.ErrLog(ctx, defined.ParamsErr, nil)
> 	}
> 	fileName := fmt.Sprintf("%s/%s%s", time.Now().Format(utils.DateFolder), utils.GenerateToken(32), ext)
> 	res, err := tencent.PresignPutObject(Config.CosUpload, fileName, 300)
> 	if err != nil {
> 		return code.ErrLog(ctx, defined.FailedErr, err)
> 	}
> 	return c.JSON(http.StatusOK, code.Success(ctx, res))
> }
> 
> ```