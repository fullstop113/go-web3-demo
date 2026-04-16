# go-web3-demo

一个用于练习 Go HTTP 服务工程化的示例项目，当前包含最小可运行的 Gin + GORM（SQLite）骨架。

## 快速开始

```bash
go run .
```

健康检查：

```bash
curl http://localhost:8080/healthz
```

## 学习指南

- [第一周周末任务执行指南（Go HTTP 服务开发）](file:///Users/liucheng/project/go-web3-demo/docs/weekend-http-service-guide.md)

<br />

<br />

<br />

- 新增文章控制器 5 个接口（增删改查 + 分页 + 作者权限校验）： article.go
- 新增中间件 RequestID ： request\_id.go
- 新增中间件 RequestLogger ： logger.go
- 新增中间件 Recovery ： recovery.go
- 新增中间件 RateLimit ： rate\_limit.go
- 路由接入全局中间件，并增加文章路由组（公开读、登录后写）： router.go
- 补充响应码常量（ forbidden/not\_found/rate\_limit ）： response.go

