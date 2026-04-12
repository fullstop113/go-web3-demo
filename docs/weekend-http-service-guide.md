# 第一周周末任务执行指南（Go HTTP 服务开发）

目标：在两天内把一个“可运行、可扩展、可持续迭代”的 Go HTTP 服务骨架搭起来。你将得到：

- RESTful API 路由分组与版本化（/api/v1）
- JWT 身份认证（Bearer Token）
- GORM 数据库操作（SQLite 起步，可切 MySQL）
- 中间件体系（日志、跨域、限流、统一错误）
- 参数校验
- 统一响应格式

本仓库 module 为：

```text
github.com/fullstop113/go-web3-demo
```

后续所有 import 路径都以此为前缀。

## 当前仓库基线（你已经拥有的最小骨架）

- main.go：启动服务、初始化 DB、初始化路由
- model/model.go：使用 GORM 连接 SQLite（app.db）
- router/router.go：Gin 路由（已有 /healthz）

建议保持这个“最小可运行基线”一直能跑通，然后在其上逐步叠加功能。

## 目标工程结构（在现有基础上扩展）

```text
go-web3-demo/
  main.go
  go.mod
  config/
  controller/
  dao/
  docs/
  middleware/
  model/
  router/
  service/
  utils/
```

你可以按功能逐步添加目录；不需要一开始把所有目录都建齐。

## Day 1：先把“认证闭环”跑通（注册、登录、鉴权）

### 1）依赖确认

本项目已经在 go.mod 中包含 Gin、GORM、SQLite、JWT、validator、rate 等依赖。若你在别的环境从零开始，可按需安装：

```bash
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/sqlite
go get github.com/golang-jwt/jwt/v5
go get github.com/go-playground/validator/v10
go get golang.org/x/time/rate
go get golang.org/x/crypto/bcrypt
```

### 2）数据模型：不要把请求体直接绑定到数据库模型

建议拆分三类结构体：

- Request：只用于接收入参（有 validate/json tag）
- Model：只用于持久化（字段更“数据库化”，敏感字段对外隐藏）
- Response：只用于返回（避免把 DB 字段与敏感字段带出去）

示例（放在 model/user.go 一类文件里）：

```go
package model

import "gorm.io/gorm"

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CreatedAt    int64          `json:"-"`
	UpdatedAt    int64          `json:"-"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Username     string         `gorm:"unique;not null" json:"username"`
	Email        string         `gorm:"unique;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
}
```

Request 示例（放在 controller 或 dto 包里都可以）：

```go
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6,max=72"`
	Email    string `json:"email" validate:"required,email"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
```

说明：

- 结构体 tag 必须写在反引号里
- 密码字段不要落库原文，永远存 PasswordHash
- Password 建议限制 max=72（bcrypt 的输入上限）

### 3）JWT：最小可用实现（并把 secret 放到环境变量）

建议在 utils/jwt.go 中实现：

- GenerateToken(userID, username) -> string
- ParseToken(tokenString) -> claims

Claims 示例：

```go
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
```

过期时间写法：

```go
jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour))
```

secret 读取策略（推荐）：

- 启动时从 JWT_SECRET 读取
- 为空就直接退出并打印清晰错误信息

### 4）统一响应：让 handler 不再重复写 JSON

约定统一响应结构：

```go
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}
```

并提供两个辅助函数：

- OK(c, data)
- Fail(c, httpStatus, code, msg)

目标是把“错误码、返回结构”集中到一个地方，避免散落在每个 handler。

### 5）鉴权中间件：只负责解析 + 注入上下文

约定 Header：`Authorization: Bearer <token>`，中间件职责：

- 缺 header / 格式错误 -> 401
- token 无效/过期 -> 401
- 正常则把 user_id / username 放入 gin.Context

建议 c.Set 的 key 固定为：

- user_id（uint）
- username（string）

### 6）路由：先只做 v1 分组 + public/private

建议路由形态：

- /healthz（保持你现有的）
- /api/v1/register（public）
- /api/v1/login（public）
- /api/v1/user/me（private）

private group 统一挂 JWT 中间件。

## Day 2：文章 CRUD + 中间件体系

### 1）文章模型与 CRUD

建议最小字段：

- title
- content
- author_id（与 user 关联）

接口建议：

- POST /api/v1/articles
- GET /api/v1/articles
- GET /api/v1/articles/:id
- PUT /api/v1/articles/:id
- DELETE /api/v1/articles/:id

实现顺序建议：

1. Create
2. List（分页）
3. Detail
4. Update（只允许作者）
5. Delete（只允许作者）

### 2）中间件清单（按“先稳定再扩展”）

建议优先级：

1. Recovery：统一捕获 panic，返回统一错误响应
2. Logger：记录请求方法、路径、耗时、状态码（避免输出敏感信息）
3. CORS：允许前端调用
4. RateLimiter：对全局或关键接口限流

注意：`gin.Default()` 自带 Logger + Recovery。如果你要自定义 Logger/Recovery，建议改为 `gin.New()` 并显式挂载。

### 3）配置：先做一个最小 config 层

建议使用环境变量配置以下内容：

- HTTP_ADDR（默认 :8080）
- DB_DSN（SQLite 默认 app.db，后续可切 MySQL DSN）
- JWT_SECRET（必填）

读取方式保持简单：启动时读 env，缺了就失败退出。

## 最小验收清单（自查）

- /healthz 返回 200 且结构稳定
- 注册成功后能登录拿到 token
- private 路由未带 token 返回 401
- private 路由带 token 能拿到当前用户信息
- 文章 CRUD 可用，且 Update/Delete 做了“只允许作者”校验
- 响应结构统一（code/msg/data）
- 关键错误不泄露内部细节（例如数据库/签名失败细节）

## 常见坑（建议你在写代码时刻意规避）

- 结构体 tag 未加反引号导致编译失败
- import 路径与 go.mod module 不一致导致编译失败
- 将请求结构体和数据库模型混用导致字段污染与敏感信息泄露
- gin.Default() 与自定义 Logger/Recovery 叠加导致重复日志
- JWT secret 写死在代码里导致泄露风险

## 进阶路线（完成后再做）

- 性能分析：pprof（可选启用、受鉴权保护）
- 可观测性：Prometheus 指标 + 结构化日志
- 数据库迁移：引入迁移工具，替代 AutoMigrate 的“隐式变更”
- 微服务化：gRPC、服务发现、配置中心

