# Gorequest From Yaml

## Overview

1. **主题:** 从 Yaml 文件读取并处理 API 请求
2. **目的:** 方便接口调用，让开发者仅关注数据层的逻辑处理，而非请求体的构建


## Layout

- The go project directory structure.
```text
.
|-- README.md
|-- configs
|   `-- requests.yml           // HTTP服务API请求配置文件
|-- go.mod
|-- main.go                    // 入口
|-- pkg
|   |-- requests               // 下载最新 go_requests 的包
|   |   |-- request.go
|   |   `-- response.go
|   `-- yaml_requests
|       `-- request.go
`-- test
    `-- unit_test              // 测试示例
        `-- demo.go
```


## Deploy using

```bash
# 修改为gomod包管理方式
go env -w GO111MODULE="on"

# 进入项目根目录
cd ./gorequest-from-yaml

# 初始化项目, 生成go.mod
go mod init ${PWD}

# 检查依赖描述
go mod tidy

# 将依赖下载至本地
go mod download

# 编写接口配置
vim ./configs/requests.yml

# 编写请求
vim ./main.go

# 编译打包
go build -o demo main.go

# 运行
./demo
```
