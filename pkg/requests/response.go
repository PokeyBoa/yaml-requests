package requests

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type Response struct {
	// http response对象
	R *http.Response

	// response状态码
	StatusCode int

	byte []byte
	text string
}

// Byte 获取Response Body为byte数组
func (resp *Response) Byte() []byte {

	var err error

	if len(resp.byte) > 0 {
		return resp.byte
	}

	var Body = resp.R.Body
	resp.byte, err = ioutil.ReadAll(Body)
	if err != nil {
		return nil
	}

	return resp.byte
}

// Text 获取Response Body为String对象
func (resp *Response) Text() string {
	if resp.byte == nil {
		resp.Byte()
	}
	resp.text = string(resp.byte)
	return resp.text
}

// SaveFile 将Response Body转换为文件
func (resp *Response) SaveFile(filename string) error {
	if resp.byte == nil {
		resp.Byte()
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(resp.byte)
	f.Sync()

	return err
}

// JsonUnmarshal 将Response Body序列化为一个struct对象
func (resp *Response) JsonUnmarshal(v interface{}) error {
	if resp.byte == nil {
		resp.Byte()
	}
	return json.Unmarshal(resp.byte, v)
}
