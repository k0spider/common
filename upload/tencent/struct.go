package tencent

type CosRequest struct {
	FileName string `json:"file_name" validate:"required"`
}

type CosResponse struct {
	Url       string `json:"url"`
	SignedUrl string `json:"signed_url"`
}
