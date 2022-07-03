package yaml_requests

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"gorequest-from-yaml/pkg/requests"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

//@author: [chenbo.mystic]
//@function: HTTP_Requests_from_yaml
//@description: 通过YAML编写HTTP API，进行数据交互

/* read yaml */
type GetBasic struct {
	Host    string                 `yaml:"host"`
	Route   string                 `yaml:"route"`
	Method  string                 `yaml:"method"`
	Params  map[string]interface{} `yaml:"params"`
	Auth    map[string]interface{} `yaml:"auth"`
	Headers map[string]string      `yaml:"headers"`
	Payload map[string]interface{} `yaml:"payload"`
}

func (s *GetBasic) httpRequest(filename string) map[string]GetBasic {
	yfile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	data := make(map[string]GetBasic)
	err2 := yaml.Unmarshal(yfile, &data)
	if err2 != nil {
		log.Fatal(err2)
	}
	return data
}

/* convert */
func dealUpper(s string) string {
	s = strings.ToUpper(strings.TrimSpace(s))
	return s
}

func deallower(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	return s
}

func jsonToMap(s string) map[string]interface{} {
	var tempMap map[string]interface{}
	err := json.Unmarshal([]byte(s), &tempMap)
	if err != nil {
		panic(err)
	}
	return tempMap
}

/* file and path */
func getRootPath() string {
	uplevel := 2
	sep := string(os.PathSeparator)
	abspath := getCurAbsPath()
	splitpath := strings.Split(abspath, sep)
	uplevel = len(splitpath) - uplevel
	rootpath := strings.Join(splitpath[:uplevel], sep)
	return rootpath
}

func getCurAbsPath() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

func memPath(args ...string) string {
	base := getRootPath()
	path := []string{base}
	if len(args) > 0 {
		path = append(path, args...)
	}
	filepath := strings.Join(path, string(os.PathSeparator))
	res := fileExist(filepath)
	if res == false {
		return ""
	}
	return filepath
}

func fileExist(path string) bool {
	fi, err := os.Lstat(path)
	if err == nil {
		return !fi.IsDir()
	}
	return !os.IsNotExist(err)
}

/* utils */
func getOsEnv(s string) string {
	s = os.Getenv(s)
	return s
}

func getMapKeys(m map[string]interface{}) []string {
	j := 0
	keys := make([]string, len(m))
	for k := range m {
		keys[j] = k
		j++
	}
	return keys
}

/* tools */
func parseYamlToMap(filename string, section string) map[string]interface{} {
	ins := GetBasic{}
	data := ins.httpRequest(filename)
	for key, value := range data {
		if key == section {
			valOfData := reflect.ValueOf(value)
			// 判断host语法
			host := valOfData.FieldByName("Host").String()
			if strings.HasSuffix(host, "/") == false {
				host += "/"
			}
			// 判断route语法
			route := valOfData.FieldByName("Route").String()
			if strings.HasPrefix(route, "/") == true {
				route = route[1:]
			}
			// 判断http请求方法
			method := valOfData.FieldByName("Method").String()
			// 处理请求参数: 转成空接口再转成map类型
			_params := valOfData.FieldByName("Params").Interface()
			params := _params.(map[string]interface{})
			// 处理auth
			_auth := valOfData.FieldByName("Auth").Interface()
			auth := _auth.(map[string]interface{})
			// 处理headers
			_headers := valOfData.FieldByName("Headers").Interface()
			headers := _headers.(map[string]string)
			// 处理payload
			_payload := valOfData.FieldByName("Payload").Interface()
			payload := _payload.(map[string]interface{})
			// 返回结果
			return map[string]interface{}{
				"host":    host,
				"route":   route,
				"method":  method,
				"params":  params,
				"auth":    auth,
				"headers": headers,
				"payload": payload,
			}
		}
	}
	return nil
}

func requestsCollect(filename string, section string, args ...string) map[string]interface{} {
	// 从yaml文件中读取相应的值
	data := parseYamlToMap(filename, section)
	// 初始化判断
	var (
		url         string
		set_params  bool
		set_auth    bool
		set_headers bool
		set_payload bool
	)
	formats := getMapKeys(data)
	for _, v := range formats {
		switch deallower(v) {
		case "params":
			set_params = true
		case "auth":
			set_auth = true
		case "headers":
			set_headers = true
		case "payload":
			set_payload = true
		}
	}
	// 解析文件相关的内容
	method := dealUpper(data["method"].(string))
	switch  method{
	// 目前支持get/post请求
	case "GET", "POST":
		var auth_container []string
		// 设置query params
		if set_params == false {
			url = strings.Join([]string{data["host"].(string), data["route"].(string)}, "")
		} else {
			// 拼接请求参数url
			queryParams := "?"
			num := 0
			for key, value := range data["params"].(map[string]interface{}) {
				// 定义占位符
				placeholder := "X" + strconv.Itoa(num)
				// 值为空的情况
				if value == nil {
					queryParams += strings.Join([]string{"&" + key, placeholder}, "=")
					num++
				}
				// 值不为空的情况（类型断言）
				if val, ok := value.(string); ok == true {
					queryParams += strings.Join([]string{"&" + key, val}, "=")
				}
			}
			queryParams = strings.Replace(queryParams, "?&", "?", 1)
			// 拼接get方法的完整url
			url = strings.Join([]string{data["host"].(string), data["route"].(string), queryParams}, "")
			// 将Xn占位符替换为真正的值
			for i := 0; i < len(args); i++ {
				url = strings.Replace(url, "X"+strconv.Itoa(i), args[i], 1)
			}
		}
		// 设置headers
		if set_headers == true {
			// 获取到headers
			for key, value := range data["headers"].(map[string]string) {
				token := getOsEnv(value)
				if token != "" {
					data["headers"].(map[string]string)[key] = token
				}
			}
		}
		// 设置payload
		if set_payload == true {
			//TODO: 设置yaml中指定的payload逻辑
		}
		// 设置auth
		if set_auth == true {
			// 获取到auth
			for key, value := range data["auth"].(map[string]interface{}) {
				if key == "type" {
					switch value {
					case "Basic Auth":
						auth_container = append(auth_container, "Basic Auth")
						fallthrough
					case "Bearer Token":
						if len(auth_container) == 0 {
							auth_container = append(auth_container, "Bearer Token")
						}
						goto LABEL_AUTH
					default:
						return nil
					}
				}
			}
		LABEL_AUTH:
			auth_container = append(auth_container)
			for key, value := range data["auth"].(map[string]interface{}) {
				if key == "content" {
					for i := 0; i < len(value.([]interface{})); i++ {
						auth_container = append(auth_container, value.([]interface{})[i].(string))
					}
				}
			}
		}
		// 返回结果
		return map[string]interface{}{
			"url":     url,
			"headers": data["headers"],
			"auth":    auth_container,
			"method": method,
		}
	default:
		return nil
	}
}

func YamlRequests(section string, payload map[string]string, args ...string) map[string]interface{} {
	// 判断yaml文件路径
	ymlfile := memPath("configs", "requests.yml")
	// 读取相关请求数据
	reqData := requestsCollect(ymlfile, section, args...)
	// 实例化Requests对象
	req := requests.Requests{}
	// 设置Headers
	header_params := reqData["headers"].(map[string]string)
	if header_params != nil {
		req.Header = &http.Header{}
		for key, value := range header_params {
			req.Header.Add(key, value)
		}
	}
	// 设置Authorization
	authlist := reqData["auth"].([]string)
	if authlist != nil {
		switch authlist[0] {
		case "Basic Auth":
			username := getOsEnv(authlist[1])
			password := getOsEnv(authlist[2])
			req.SetBasicAuth(username, password)
		case "Bearer Token":
			req.SetBearerToken(authlist[1])
		default:
			return nil
		}
	}
	// 获取Url
	url := reqData["url"].(string)
	// 获取method
	method := reqData["method"].(string)
	// 执行对应方法
	switch method {
	// 目前支持get/post请求
	case "GET":
		// 执行HTTP Get请求
		resp, err := req.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.R.Body.Close()
		// 判断状态码
		if resp.StatusCode == 200 {
			// 打印Response Body
			body := resp.Text()
			// 转成map
			result := jsonToMap(body)
			return result
		}
		return nil
	case "POST":
		// 准备POST Data
		/*
			requestData := map[string]string{
				"app_id":     "xxxx",
				"app_secret": "xxxx",
			}
		*/
		// 执行HTTP POST请求
		if strings.HasSuffix(url,"?") == true {
			url = url[:len(url)-1]
		}
		resp, err := req.PostJson(url, payload)
		if err != nil {
			panic(err)
		}
		defer resp.R.Body.Close()
		// 判断状态码
		if resp.StatusCode == 200 {
			//TODO: 返回真实数据
			return map[string]interface{}{
				"msg":  "success",
				"code": 0,
				"data": "{}",
			}
		}
		return nil
	default:
		return nil
	}
}
