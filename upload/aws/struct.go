package aws

type AwsS3Request struct {
	FileName string `json:"file_name" validate:"required"`
}

type AwsS3Response struct {
	Header map[string]string `json:"header"`
	Url    string            `json:"url"`
	Name   string            `json:"name"`
}
