package utils

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"time"
)

type HttpClient struct {
	Client  *resty.Client
	Request *resty.Request
}

type ResponseString struct {
	Str string
}

func NewHttp() *HttpClient {
	// Create a Resty Client
	client := resty.New()
	// Retries are configured per client
	client.SetTimeout(30 * time.Second)

	return &HttpClient{
		Client:  client,
		Request: client.R(),
	}
}

func (HttpClient) ResponseCheckErr(rep *resty.Response, err error) (*resty.Response, error) {
	if err != nil {
		return nil, err
	}
	if !rep.IsSuccess() {
		return nil, errors.Errorf("The request was unsuccessful httpStatus:%s, rep:%s", rep.Status(), rep.String())
	}
	if rep.String() == "{}" {
		return nil, errors.Errorf("Return to {}")
	}

	code := gjson.Parse(rep.String()).Get("code")
	if !code.Exists() || code.Int() == 0 || code.Int() == 200 {
		return rep, nil
	}
	return nil, errors.Errorf("Nonzero value code rep:%s", rep.String())
}

func (h *HttpClient) Get(url string, result interface{}) error {
	rep, err := h.ResponseCheckErr(h.Request.Get(url))
	if err != nil {
		return err
	}
	if s, ok := result.(*ResponseString); ok {
		s.Str = rep.String()
		return nil
	}

	if err = json.Unmarshal(rep.Body(), result); err != nil {
		return err
	}
	return nil
}

func (h *HttpClient) Post(url string, result interface{}) error {
	rep, err := h.ResponseCheckErr(h.Request.Post(url))
	if err != nil {
		return err
	}
	if err = json.Unmarshal(rep.Body(), result); err != nil {
		return err
	}
	return nil
}
