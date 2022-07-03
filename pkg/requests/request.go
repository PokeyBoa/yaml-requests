package requests

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

const (
	defaultTimeout = 60 * time.Second       // http请求默认超时时间
)

type Requests struct {

	// 设置是否跳过Server端证书验证
	SkipTLS bool

	// http请求默认超时时间
	// The timeout includes connection time, any redirects, and reading the response body
	// Timeout of zero means no timeout.
	Timeout time.Duration

	// Http Accept default is "application/json"
	Accept string

	// Http ContentType default is "application/json"
	ContentType string

	// Http request header, 如果有Header需求请设置
	Header *http.Header

	// Http request cookie, 如果有Cookie需求请设置
	Cookie *http.Cookie

	client  *http.Client
	request *http.Request
}

// setSkipTLS 设置跳过Server端证书验证
func (r *Requests) setSkipTLS() {
	if r.SkipTLS {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		r.client.Transport = tr
	}
}

// SetCookie 设置HTTP Request Cookie
func (r *Requests) setCookie() {
	if r.request.Header == nil {
		r.request.Header = http.Header{}
	}

	if r.Cookie != nil {
		r.request.AddCookie(r.Cookie)
	}
}

// SetHeader 设置HTTP Request Header
func (r *Requests) setHeader() {
	if r.request.Header == nil {
		r.request.Header = http.Header{}
	}
	if r.Header != nil {
		r.request.Header = *r.Header
	}

	if r.Accept == "" {
		r.request.Header.Set("Accept", "application/json")
	}

	if r.ContentType == "" {
		r.request.Header.Set("Content-Type", "application/json")
	}
}

// 设置http连接超时时间
func (r *Requests) setTimeout() {
	if r.Timeout == 0 {
		r.client.Timeout = defaultTimeout
	} else {
		r.client.Timeout = r.Timeout
	}
}

// SetBearerToken 设置Http Request Bearer Token
func (r *Requests) SetBearerToken(bearerToken string) {
	if r.Header == nil {
		r.Header = &http.Header{}
	}
	r.Header.Add("Authorization", "Bearer "+bearerToken)
}

// SetBasicAuth 设置Http Request Basic Auth
func (r *Requests) SetBasicAuth(username, password string) {
	if r.Header == nil {
		r.Header = &http.Header{}
	}
	r.Header.Set("Authorization", "Basic "+basicAuth(username, password))
}

// Get 执行http get请求
func (r *Requests) Get(url string) (*Response, error) {
	var err error
	r.request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return r.do()
}

// Post 执行http post请求, body是一个string对象
func (r *Requests) Post(url string, body string) (*Response, error) {
	var err error
	r.request, err = http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	return r.do()
}

// PostJson 执行http post请求, body是一个struct对象
func (r *Requests) PostJson(url string, body interface{}) (*Response, error) {
	byteBody, _ := json.Marshal(body)

	var err error
	r.request, err = http.NewRequest("POST", url, bytes.NewBuffer(byteBody))
	if err != nil {
		return nil, err
	}

	return r.do()
}

func (r *Requests) Put(url string, body string) (*Response, error) {
	var err error
	r.request, err = http.NewRequest("PUT", url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	return r.do()
}

func (r *Requests) PutJson(url string, body interface{}) (*Response, error) {
	byteBody, _ := json.Marshal(body)

	var err error
	r.request, err = http.NewRequest("PUT", url, bytes.NewBuffer(byteBody))
	if err != nil {
		return nil, err
	}

	return r.do()
}

func (r *Requests) Patch(url string, body string) (*Response, error) {
	var err error
	r.request, err = http.NewRequest("PATCH", url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	return r.do()
}

func (r *Requests) PatchJson(url string, body interface{}) (*Response, error) {
	byteBody, _ := json.Marshal(body)

	var err error
	r.request, err = http.NewRequest("PATCH", url, bytes.NewBuffer(byteBody))
	if err != nil {
		return nil, err
	}

	return r.do()
}

func (r *Requests) Delete(url string) (*Response, error) {
	var err error
	r.request, err = http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	return r.do()
}

// Do send an HTTP request and returns an HTTP response
func (r *Requests) do() (*Response, error) {

	// 实例化http client
	r.client = &http.Client{}

	// 参数加载
	r.setSkipTLS()
	r.setTimeout()
	r.setHeader()
	r.setCookie()

	res, err := r.client.Do(r.request)
	if err != nil {
		return nil, err
	}

	resp := &Response{}
	resp.R = res
	resp.StatusCode = resp.R.StatusCode
	resp.Byte()
	return resp, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
