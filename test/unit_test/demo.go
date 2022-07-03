package unit_test

import (
	"fmt"
	"gorequest-from-yaml/pkg/yaml_requests"
)

// API接口测试示例
func GetReq() {
	// 百度百科搜索 '关键字'
	response := yaml_requests.YamlRequests("baidu_baike", nil, "以父之名")
	// 返回请求数据
	for k, v := range response {
		fmt.Printf("[key]: %v, [value]: %v\n", k, v)
	}
}

func PostReq() {
	// 准备POST Data
	payload := map[string]string{
		"username": "xxxxxxxx",
		"password": "xxxxxxxx",
	}
	// 发送post请求
	response := yaml_requests.YamlRequests("post_auth", payload)
	// 返回请求数据
	fmt.Println(response)
}
