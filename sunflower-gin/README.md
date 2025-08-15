# gin-base-layout

基于gin的项目布局基础代码，采用最简单最朴素最好理解的目录结构，没有使用依赖注入。

脚手架工具 👉 [iaa](https://github.com/q1mi/iaa)

## 目录结构

```bash
├── api
│   ├── api.go
│   ├── calc
│   └── code.go
├── cmd
│   ├── gen
│   └── server
├── config
│   └── config.yaml
├── deploy
├── docs
├── go.mod
├── go.sum
├── internal
│   ├── conf
│   ├── dao
│   ├── handler
│   ├── middleware
│   ├── model
│   ├── server
│   ├── service
│   └── task
├── LICENSE
├── log
│   └── server.log
├── Makefile
├── pkg
│   ├── jwt
│   ├── logging
│   └── snowflake
├── scripts
└── test
```

## 快速开始

1. 修改配置文件 `config/config.yaml`
2. 运行服务
```bash
go run cmd/server/main.go
```