# go-admin

面向业务开发的 Go 管理后台框架，配合 [go-admin-template](https://github.com/Gary-Yez/go-admin-template) 快速搭建后台项目。

- **快速生成业务代码**：生成 Go 模块、Vue 页面、API 调用和菜单，按需选择增删改能力。
- **内置权限管理**：管理员多角色、角色切换、菜单授权、API 自动登记及个人 API 密钥。
- **类型明确的业务配置**：通过字段读取单项配置，支持数据库持久化和动态初始值。
- **统一基础能力**：复用数据库、内存或 Redis 缓存，以及定时任务调度。

## 从哪里开始

开发业务项目，先使用模板仓库完成前后端启动，再按本文调用框架提供的方法。

```text
go-admin/
├── admin.go       # 服务启动和模块注册
├── module.go      # 模块接口
├── menu.go        # 菜单声明
├── configs.go     # 内置配置和配置项类型
├── services.go    # 数据库、缓存、调度器入口
├── request/       # 请求参数与当前身份
├── response/      # 统一响应
└── internal/      # 框架内部实现，业务无需直接引用
```

业务主要导入：

```go
import (
    admin "github.com/Gary-Yez/go-admin"
    "github.com/Gary-Yez/go-admin/request"
    "github.com/Gary-Yez/go-admin/response"
)
```

## 公开方法速查

| 方法或类型 | 用途 |
| --- | --- |
| `admin.Run()` | 读取配置并启动后台服务 |
| `admin.Register(key, module)` | 注册业务模块，返回错误 |
| `admin.MustRegister(key, module)` | 注册模块，错误时 panic |
| `admin.Module` | 业务模块需要实现的接口 |
| `admin.MenuDefinition` | 声明模块默认菜单 |
| `admin.DB()` | 获取 GORM 数据库连接 |
| `admin.Cache()` | 获取统一缓存接口 |
| `admin.Scheduler()` | 获取任务调度器 |
| `admin.ConfigItem[T]` | 带具体值类型的配置项 |
| `admin.ConfigDefault(value)` | 在配置 Init 中构造动态初始值 |
| `admin.RegisterConfig(&config)` | 绑定并注册配置结构，模板已自动调用 |
| `request.GetAuthUser(ctx)` | 获取当前用户和角色 |
| `request.GetReqList(ctx)` | 绑定分页、筛选、排序参数 |
| `request.GetReqIds(ctx)` | 绑定并去重批量 ID |
| `request.GetReq(ctx)` | 绑定单个 ID |
| `response.Success / List / Error` | 返回统一结果 |

## 启动服务

当前项目使用 Go 1.25.5。准备 MySQL 或 PostgreSQL 数据库后，在业务项目中调用：

```go
func main() {
    if err := admin.Run(); err != nil {
        panic(err)
    }
}
```

默认读取工作目录下的 `config.yaml`。文件不存在时会创建默认文件并退出，填写连接信息后重新启动。

```go
admin.Run("config.dev.yaml") // 也可以在代码里指定文件
```

```sh
go run . -c config.dev.yaml
go run . --config config.prod.yaml --server.port 9000
```

命令行配置文件优先于代码参数。YAML 填写数据库、Redis、监听端口等环境参数；业务配置在后台维护。完整 YAML 和前后端启动步骤见模板 README。

数据库统一使用 `database` 配置，业务仍通过 `admin.DB()` 查询：

```yaml
database:
  driver: mysql # mysql / postgres
  host: 127.0.0.1
  port: "3306" # PostgreSQL 使用 5432；省略时按 driver 选择
  name: go_admin
  username: your_user
  password: your_password
  sslmode: disable # 仅 PostgreSQL 使用，可按服务端要求配置 require/verify-full
```

已有配置需将 `mysql` 节点改为 `database`，其中 `database` 字段改名为 `name`，并添加 `driver`。这只切换连接，不会将原数据库的数据搬迁到另一种数据库。环境变量使用 `MYAPP_DATABASE_DRIVER`、`MYAPP_DATABASE_HOST`、`MYAPP_DATABASE_PORT`、`MYAPP_DATABASE_NAME` 等。

缓存、分布式锁和权限通知的命名空间包含数据库类型、主机、端口及库名；同一部署的实例应使用一致的连接标识。升级后旧命名空间缓存不再读取，由现有过期机制清理。

## 注册业务模块

模块需要实现四个方法：

```go
type Module interface {
    Name() string
    Initialize() error
    AdminRouter(*gin.RouterGroup)
    PublicRouter(*gin.RouterGroup)
}
```

| 方法 | 应该放什么 |
| --- | --- |
| `Name()` | 模块名称 |
| `Initialize()` | 表迁移、任务处理函数注册等启动操作 |
| `AdminRouter()` | 需要登录和角色权限的接口 |
| `PublicRouter()` | 公开接口 |

在模板的 `modules/enter.go` 注册，业务 import 路径按项目 module 调整：

```go
package modules

import (
    admin "github.com/Gary-Yez/go-admin"
    "example.com/app/modules/order"
)

func init() {
    admin.MustRegister("order", new(order.Mounter))
}
```

主程序匿名导入 modules 即可自动注册。包 `init()` 只注册；数据库和缓存相关操作放在 `Initialize()`，此时框架依赖已经就绪。

例如模块 Key 为 order，在 AdminRouter 中注册 `group.POST("list", Controller.List)`，默认完整地址就是 `/api/order/list`。

### 声明默认菜单

模块可额外实现 `Menus()`：

```go
func (*Mounter) Menus() []admin.MenuDefinition {
    return []admin.MenuDefinition{{
        Key: "order", Name: "订单管理",
        Icon: "iconoir:page", Path: "order",
        Component: "../views/order/index.vue", Sort: 10,
    }}
}
```

父菜单填写 `ParentKey`。框架按 Key 补建菜单，保留后台对已有菜单的修改。生成器可以自动生成这个方法。普通角色仍需要分别授权菜单和 API。

## 使用业务配置

模板已经生成配置项并自动注册，业务直接导入本地 settings 包：

```go
name, err := settings.BaseConfig.SiteName.Get() // string, error
minutes := settings.BaseConfig.JwtExpireMinutes.MustGet() // int
text := settings.Test.MustGet() // 模板当前的示例 string 配置
```

内置配置通过 `settings.BaseConfig` 访问，用户配置直接通过 `settings.字段名` 访问。每次只读取对应 Key；返回类型由字段定义确定，不需要类型断言。

| 方法 | 行为 |
| --- | --- |
| `Get()` | 返回具体类型的值和错误 |
| `MustGet()` | 返回具体类型的值，读取失败时 panic |

配置支持 string、bool、int、float64、[]string。读取优先缓存，未命中时查询数据库；在配置管理页面修改实际值。

配置定义页面支持拖拽业务分组和组内配置项，预览确认后按字段声明顺序保存到代码。配置管理按保存顺序展示；换数据库部署时也会按代码恢复顺序。内置配置及其分组保持 BaseConfig 声明顺序。

### 动态初始值

在开发工具添加 `order.timeout`（int）后，可在模板的 `settings/init.go` 中填写：

```go
package settings

import admin "github.com/Gary-Yez/go-admin"

func (c *config) Init() error {
    c.OrderTimeout = admin.ConfigDefault(60)
    return nil
}
```

固定默认值直接在配置定义页面填写。动态初始值由 Init 构造，只补建缺失配置，不覆盖数据库已有值。框架内置配置自动初始化，用户不需要手动调用。此方法中不要通过 Get/MustGet 读取配置。

JWT 签名密钥位于启动配置 `jwt.secret`，生成配置文件时自动随机填写，也可通过 `MYAPP_JWT_SECRET` 环境变量覆盖。密钥至少 32 字节，多实例须配置一致；修改后重启，旧 JWT 失效。真实密钥不要提交到仓库。登录有效期保留在配置管理，默认 10080 分钟，仅影响新签发的令牌。

## 数据库与请求参数

### 数据库

`admin.DB()` 返回 `*gorm.DB`，使用 GORM 原有方法操作：

```go
// 模块 Initialize 中迁移业务表。
return admin.DB().AutoMigrate(&Order{})
```

### 当前身份

在控制器中：

```go
authUser, err := request.GetAuthUser(ctx)
if err != nil {
    response.Error(ctx, err, 401)
    return
}
// authUser.UserId：当前用户 ID
// authUser.RoleId：本次请求的角色 ID
```

### 列表查询

以下片段放在使用 `Order` 模型的列表控制器中：

```go
req, err := request.GetReqList(ctx)
if err != nil {
    response.Error(ctx, err)
    return
}

db := req.WithFilter(admin.DB().Model(&Order{}), []string{"status"})
var total int64
if err := db.Count(&total).Error; err != nil {
    response.Error(ctx, err)
    return
}
var list []Order
err = req.WithPagination(req.WithSort(db, []string{"id"})).Find(&list).Error
if err != nil {
    response.Error(ctx, err)
    return
}
response.List(ctx, list, total)
```

`WithFilter` 和 `WithSort` 的字段列表是白名单。非法字段或条件会返回查询错误，不能忽略 GORM Error。分页包含 page、limit，limit 最多 100；filters 支持比较、like、in、between，sorts 使用 field 和 order。

批量操作使用 `request.GetReqIds(ctx)`，通过 `req.WithQuery(admin.DB())` 生成 ID 查询条件。ID 要求非空正整数，自动去重，不限制数量。单个 ID 使用 `request.GetReq(ctx)`。

### 返回结果

```go
response.Success(ctx)                         // 成功
response.Success(ctx, "保存成功")             // 自定义 message
response.Success(ctx, gin.H{"id": 1})        // 返回 data，需要导入 gin
response.List(ctx, list, total)                // data: {list, total}
response.Error(ctx, err)                       // 错误
response.Error(ctx, "没有权限", 403)           // 指定业务错误码
```

响应包含 code、message 和可选 data。当前统一响应使用 HTTP 200，调用方应检查 JSON `code`，成功值为 200。

## 使用缓存

通过 `admin.Cache()` 调用，无需区分内存和 Redis：

键只需填写 `模块:具体键`。框架根据数据库地址、端口和库名计算命名空间，普通缓存实际为 `go-admin:<数据库哈希>:cache:<业务键>`，锁为 `go-admin:<数据库哈希>:lock:<业务键>`。同一系统的多实例应使用一致的数据库连接地址。通过 `Client()` 直接调用原始客户端时，不会自动添加这些前缀。

```go
// 片段需要导入 time、errors，并处理各操作的 err。
err := admin.Cache().Set("order:last_id", "123", time.Minute)
value, err := admin.Cache().Get("order:last_id")
if errors.Is(err, admin.ErrCacheNotFound) {
    // 缓存不存在，按业务需要回查数据库。
}
```

| 方法 | 用途 |
| --- | --- |
| `Set(key, value, ttl)` / `Get(key)` | 写入值、读取字符串 |
| `SetJSON(key, value, ttl)` / `GetJSON(key, &value)` | 结构体序列化读写 |
| `GetInt(key)` / `GetBool(key)` | 读取整数或布尔值 |
| `Exists(key)` / `Del(key)` | 检查或删除缓存 |
| `Lock(key, expiration, extendInterval)` | 获取锁，使用后调用返回对象的 Unlock |

TTL 为 0 表示不过期。单实例可使用内存缓存，多实例配置共享 Redis。推荐处理返回错误；缓存的 MustGet 系列会忽略错误，与配置项 MustGet 的 panic 行为不同。

## 注册定时任务

在模块 `Initialize()` 中注册处理函数，随后在后台创建任务并选择它。以下为方法体示例，需要导入 context：

```go
return admin.Scheduler().RegisterHandler("order.cleanup", &admin.HandlerOption{
    Name: "清理过期订单",
    Handler: func(ctx context.Context, params []byte) error {
        if err := ctx.Err(); err != nil {
            return err
        }
        // 执行清理逻辑；有参数时将 params 按 JSON 解码。
        return nil
    },
})
```

`HandlerOption.Params` 可声明参数，使用 `admin.HandlerParams`、`admin.HandlerParam`，类型有 `admin.StringParams`、`admin.IntParams`、`admin.BoolParams`。处理函数 Key 必须唯一。

后台 Cron 使用五段格式，例如 `*/5 * * * *` 每五分钟执行。默认 UTC，也可填写 `CRON_TZ=Asia/Shanghai 0 2 * * *`。任务执行结果在后台查看，业务处理函数应响应 context 并保证业务幂等。

## 配置全局 Gin 服务

在 `main()` 中、`admin.Run()` 前调用 `admin.ConfigureEngine(func(*gin.Engine) error)`。可以注册多个回调，按顺序在路由注册前执行；回调报错则停止启动。服务启动后不能再注册。可用于设置可信代理、添加全局中间件等。

框架不再自动安装 CORS；模板 `server/main.go` 使用此回调配置 CORS（包含 Authorization 请求头），开发者可以在那里调整允许的来源。回调直接注册的路由不会自动获得框架鉴权，也不参与 API 同步，业务接口应继续通过模块的 `AdminRouter` / `PublicRouter` 注册。

## 节点监控

“系统运维 → 节点监控”展示各实例的资源和 Go Runtime 数据。每 5 秒采集，超过 20 秒没有上报标记离线，Redis 中保留最后快照 10 分钟。单实例使用内存，多实例按现有部署约定共用 Redis 及数据库命名空间；节点目录使用 Redis Hash + Sorted Set 原子更新，不扫描全库、不写 MySQL。

`server.node_name`（或 `MYAPP_SERVER_NODE_NAME`）设置显示名称，留空使用主机名；实例 ID 每次启动随机生成，避免同一主机多进程相互覆盖。重启后旧实例暂时显示离线，10 分钟后移除。节点时钟应同步；Redis 在线状态以 Redis 时间计算。

机器指标是操作系统可见资源，不代表容器 CPU/内存配额；磁盘展示运行目录所在卷。进程 CPU 按单核 100% 计算，可超过 100%。采集不支持或失败显示不可用，首次 CPU 和 GC 增量需要等待下一次采样。Go 堆内存不等同于进程 RSS。

监控接口 `GET /sys_monitor/list` 属于受保护路由，需要单独分配角色权限；不提供 pprof、环境变量、密钥或连接凭据。Redis 故障时不会伪装成单机列表，上报会持续重试，页面展示读取失败。页面支持节点详情和当前页面内的短趋势，离开页面不再轮询，不保存历史监控数据。
