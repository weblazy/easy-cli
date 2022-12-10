# easy-cli

# Features
- APIs ：协议通信以 HTTP/gRPC 为基础，通过 Protobuf 进行定义；
- Errors ：通过 Protobuf 的 Enum 作为错误码定义，以及工具生成判定接口；
- Metadata ：在协议通信 HTTP/gRPC 中，通过 Middleware 规范化服务元信息传递；
- Config ：支持多数据源方式，进行配置合并铺平，通过 Atomic 方式支持动态配置；
- Logger ：标准日志接口，可方便集成三方 log 库，并可通过 fluentd 收集日志；
- Metrics ：统一指标接口，可以实现各种指标系统，默认集成 Prometheus；
- Tracing ：遵循 OpenTelemetry 规范定义，以实现微服务链路追踪；
- Encoding ：支持 Accept 和 Content-Type 进行自动选择内容编码；
- Transport ：通用的 HTTP /gRPC 传输层，实现统一的 Middleware 插件支持；
- Registry ：实现统一注册中心接口，可插件化对接各种注册中心；
Validation: 通过Protobuf统一定义校验规则，并同时适用于HTTP/gRPC服务.
- SwaggerAPI: 通过集成第三方Swagger插件 能够自动生成Swagger API json并启动一个内置的Swagger UI服务.
- 集成postman生成
- 集成yapi生成
- gorm orm文件生成
- http服务生成
- grpc服务生成
- http mock
- grpc mock
- 监控面板
- 告警
# 目录结构
.
├── Dockerfile
├── README.md
├── common
│   └── common.go
├── conf
│   ├── base.go
│   ├── const.go
│   └── local.go
├── cronjobs
│   ├── config
│   │   └── config.go
│   ├── cronjobs.go
│   └── handler
│       └── sync_user.go
├── easy.yaml
├── go.mod
├── go.sum
├── grpcs
│   └── order_rpc
│       ├── config
│       │   └── config.go
│       ├── handler
│       │   └── user.go
│       ├── logic
│       │   └── user_logic
│       │       └── get_user_info.go
│       ├── order_rpc.go
│       └── proto
│           └── user
│               ├── user.pb.go
│               ├── user.proto
│               └── user_grpc.pb.go
├── https
│   └── order_http
│       ├── config
│       │   └── config.go
│       ├── def
│       │   └── def.go
│       ├── errcode
│       │   └── errcode.go
│       ├── handler
│       │   └── user.go
│       ├── logic
│       │   └── user
│       │       └── get_user_info.go
│       ├── order_http.go
│       └── routes
│           ├── routes.go
│           └── user_interceptor.go
├── jobs
│   ├── config
│   │   └── config.go
│   ├── handler
│   │   └── init_user.go
│   └── jobs.go
├── main.go
├── model
│   ├── app
│   │   ├── mysql_client.go
│   │   └── user.go
├── orm.yaml
├── pkg
└── postman.json
