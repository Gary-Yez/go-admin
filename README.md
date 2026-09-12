# go-admin

基于 Gin、GORM 和 Casbin 的 Go 后台开发框架。业务项目通过模块注册扩展接口、菜单、数据模型和任务，系统管理能力由框架提供。

- 内置管理员、多角色、菜单与 API 权限、个人资料和 API 密钥管理。
- 支持 MySQL / PostgreSQL；单实例使用内存缓存，多实例共享 Redis。
- 提供代码生成、配置定义、配置管理、登录日志、计划任务和节点监控。
- 业务配置持久化到数据库，通过带类型的配置项按需读取缓存。

需要完整可运行项目时，从 [go-admin-template](https://github.com/Gary-Yez/go-admin-template) 开始；公共前端接入见 [go-admin-web](https://github.com/Gary-Yez/go-admin-web)。

## 安装与最小启动

要求 Go 1.25.5 或更新的兼容版本，以及可连接的 MySQL 或 PostgreSQL 数据库。框架会建表，数据库本身需要先创建。

```sh
mkdir my-admin
cd my-admin
go mod init example.com/my-admin
go get github.com/Gary-Yez/go-admin@latest
```

创建 `main.go`：

```go
package main

import (
    "log"

    admin "github.com/Gary-Yez/go-admin"
)

func main() {
    if err := admin.Run(); err != nil {
        log.Fatal(err)
    }
}
```

运行 `go run .`。工作目录没有 `config.yaml` 时，框架生成包含随机 JWT 签名密钥的配置文件并退出；填写数据库连接后重新运行。

默认监听 `0.0.0.0:8080`，API 前缀为 `/api`，管理端静态资源挂载在 `/admin`。框架包不自带编译后的前端，需要把前端产物放在运行目录的 `dist/`。

配置文件中：

- `database.driver` 为 `mysql` 或 `postgres`，对应端口通常为 3306 或 5432；PostgreSQL 还可配置 `sslmode`。
- `redis.host` 留空使用内存缓存；填写后使用 Redis。
- `server.dev` 默认关闭。开发工具需要开启它；开发接口可操作本地源码，生产环境应保持关闭。
- `jwt.secret` 至少 32 字节，修改后重启生效。登录有效期属于数据库业务配置。

## 启动接口与调用顺序

| 公开接口 | 用途 | 调用时机 |
| --- | --- | --- |
| `Register(key, module) error` | 注册业务模块，拒绝空 Key、重复 Key 和无效模块 | `Run` 前 |
| `MustRegister(key, module)` | 注册失败时 panic | `Run` 前，通常在业务模块入口 |
| `RegisterConfig(&config) error` | 注册唯一的业务配置结构并绑定字段 Key | `Run` 前 |
| `ConfigureEngine(func(*gin.Engine) error) error` | 添加全局 Gin 配置回调，可注册多个 | `Run` 前 |
| `Run(configPaths ...string) error` | 初始化依赖、模块和路由，然后阻塞提供 HTTP 服务 | main 中调用一次 |
| `DB()` / `Cache()` / `Scheduler()` | 获取框架共享服务 | 模块 `Initialize` 或更晚 |

`Run()` 默认读取工作目录的 `config.yaml`；`Run("custom.yaml")` 指定默认路径。命令行 `--config` / `-c` 优先于函数参数：

```sh
go run . --config config.production.yaml
go run . -c config.production.yaml --server.port 8081
```

配置路径不会改变工作目录，`dist/` 仍相对于进程工作目录。一个进程只调用一次 `Run`，启动失败后也不要原地再次调用。

实际启动顺序：

1. 读取环境配置，初始化数据库、缓存、计划任务服务。
2. 注册系统模块；系统模块排在业务模块前。
3. 初始化业务配置：解析默认值、执行初始化、补建数据库中缺失的配置。
4. 按顺序调用模块 `Initialize()`；系统模块先执行，业务模块按注册顺序执行。
5. 收集模块菜单，补充缺失菜单。
6. 创建 `gin.Default()`，依次执行 `ConfigureEngine` 回调。
7. 挂载静态资源，注册所有 `AdminRouter`，记录受保护路由，再注册 `PublicRouter`。
8. 同步新增 API，启动节点上报、计划任务调度和权限重新同步任务。
9. 启动 HTTP 服务并阻塞；正常返回时执行框架已注册的后台服务清理。

Go 的 `init()` 负责“声明和注册”，模块的 `Initialize()` 负责“依赖已经就绪后的初始化”。不要在包级变量或 `init()` 中调用数据库、缓存或读取配置值。框架当前没有公开的模块关闭回调，也没有通用登录事件钩子。

## 编写业务模块

模块实现 `admin.Module`：

```go
type Module interface {
    Name() string
    Initialize() error
    AdminRouter(*gin.RouterGroup)
    PublicRouter(*gin.RouterGroup)
}
```

例如创建 `modules/product/enter.go`：

```go
package product

import (
    admin "github.com/Gary-Yez/go-admin"
    "github.com/Gary-Yez/go-admin/response"
    "github.com/gin-gonic/gin"
)

type Mounter struct{}

func (*Mounter) Name() string { return "商品管理" }
func (*Mounter) Initialize() error { return nil }

func (*Mounter) AdminRouter(group *gin.RouterGroup) {
    group.GET("ping", func(ctx *gin.Context) {
        response.Success(ctx, gin.H{"message": "商品模块已就绪"})
    })
}

func (*Mounter) PublicRouter(group *gin.RouterGroup) {}

func Register() {
    admin.MustRegister("product", &Mounter{})
}
```

在 main 中 `Run()` 前调用 `product.Register()`，或者由业务 `modules` 包的 `init()` 统一注册，再让 main 空白导入该包。

以上接口路径为 `/api/product/ping`：

- `AdminRouter` 自动经过 JWT / API Token 身份校验和 Casbin 权限校验，并参与 API 同步。
- `PublicRouter` 不自动附加上述校验，也不参与受保护 API 同步。需要登录的业务接口通常放在 `AdminRouter`。
- `Name()` 用于模块展示和 API 分组，填写面向用户的中文名称。
- 新增路由进入 API 管理不等于普通角色自动获得权限；仍需分配权限。
- API 权限只控制接口访问。租户、数据归属等行级权限需要业务查询自行限制。

### 声明菜单

模块额外实现 `admin.MenuProvider` 即可：

```go
func (*Mounter) Menus() []admin.MenuDefinition {
    return []admin.MenuDefinition{{
        Key: "product", Name: "商品管理", Path: "product",
        Icon: "iconoir:box", Component: "../views/product/index.vue", Sort: 10,
    }}
}
```

`MenuDefinition` 字段为 `Key、Name、ParentKey、Icon、Path、Component、Sort、Hidden`。菜单 Key 全局唯一；父子菜单通过 `ParentKey` 关联。业务组件 Key 对应前端注册的 Vue 文件，不是浏览器 URL。

启动只补齐缺失菜单，不覆盖后台已经修改的名称、图标和排序。只要代码中保留定义，删除数据库菜单后，下次启动会重新补齐。菜单可见性与后端接口权限分别配置。

## 数据库、请求与响应

`admin.DB()` 返回共享 `*gorm.DB`。在 `Initialize()` 中建业务表，在请求中使用 `WithContext(ctx.Request.Context())`。事务内的所有数据库操作都使用回调传入的 `tx`，不要重新调用 `admin.DB()`。

`request` 包提供：

| 接口或类型 | 作用 |
| --- | --- |
| `GetReq(ctx)` / `Req.WithQuery(db)` | 绑定单个 `id`，添加主键条件 |
| `GetReqIds(ctx)` / `ReqIds.WithQuery(db)` | 绑定 `ids`，校验正整数并去重，添加批量条件 |
| `GetReqList(ctx)` / `ReqList.Validate()` | 绑定并校验分页、筛选、排序 |
| `ReqList.WithFilter(db, allowFields)` | 仅允许白名单字段筛选 |
| `ReqList.WithSort(db, allowFields)` | 仅允许白名单字段排序 |
| `ReqList.WithPagination(db)` | Page > 0 时分页；Limit 默认 10，最大 100 |
| `DefaultSortFields(extra ...string)` | 返回 id、created_at、updated_at，再追加业务字段 |
| `GetAuthUser(ctx)` | 获取已认证身份，可读取 UserId、RoleId；失败返回错误 |
| `SetAuthUser(ctx, user)` | 写入上下文身份；本身不验证令牌，不应用于绕过鉴权 |

批量 ID 不限制数量，但非空且均需大于零。Page 为零不分页；列表接口如需始终分页，应主动设为 1。

列表控制器示例（`Product` 为业务模型）：

```go
func List(ctx *gin.Context) {
    req, err := request.GetReqList(ctx)
    if err != nil {
        response.Error(ctx, err, 400)
        return
    }
    if req.Page == 0 { req.Page = 1 }
    db := req.WithFilter(admin.DB().WithContext(ctx.Request.Context()).
        Model(&Product{}), []string{"name"})
    var total int64
    if err := db.Count(&total).Error; err != nil {
        response.Error(ctx, err)
        return
    }
    var rows []Product
    db = req.WithSort(db, request.DefaultSortFields())
    if err := req.WithPagination(db).Order("id DESC").Find(&rows).Error; err != nil {
        response.Error(ctx, err)
        return
    }
    response.List(ctx, rows, total)
}
```

筛选支持 `=、!=、>、<、>=、<=、like、in、between`；排序使用 `asc / desc`。字段白名单必须由服务端定义，不要接收客户端提供的白名单。默认排序列表也不会自动为模型添加时间字段。

`response` 包的响应约定：

- `Success(ctx, data)`：`{code:200,message:"success",data:...}`。如果 data 是字符串，它作为 message；需要返回字符串数据时用对象包装。
- `List(ctx, rows, total)`：`{code:200,message:"success",data:{list:...,total:...}}`。
- `Error(ctx, err, code...)`：默认业务 code 为 500，可指定 400 / 401 / 403 等；**HTTP 状态仍为 200**，客户端需要检查 JSON code。
- `Success`、`Error` 会 Abort Gin 后续处理链，但不会替当前 Go 函数执行 return；错误分支仍要显式 return。

## 带类型的业务配置

环境连接信息放在 `config.yaml`；网站名称、登录有效期等业务值由数据库保存，配置管理页面负责修改。对业务公开的是读取接口，不是内部配置存储包。

```go
package settings

import admin "github.com/Gary-Yez/go-admin"

type config struct {
    admin.BaseConfig
    OrderTimeout admin.ConfigItem[int] `config:"order.timeout" label:"订单超时" group:"订单设置" default:"30" description:"订单未付款超时分钟数"`
}

var current config
var (
    BaseConfig = &current.BaseConfig
    OrderTimeout = &current.OrderTimeout
)

func init() {
    if err := admin.RegisterConfig(&current); err != nil { panic(err) }
}

func (c *config) Init() error {
    // 可用环境变量或计算结果填充动态初始值。
    c.OrderTimeout = admin.ConfigDefault(60)
    return nil
}
```

业务使用 `settings.OrderTimeout.Get()` 获取 `(int, error)`，或 `settings.BaseConfig.SiteName.Get()` 获取 `(string, error)`。`MustGet()` 返回同样的具体类型，但错误时 panic；请求处理优先使用 Get 并处理错误。

支持 `string、bool、int、float64、[]string`。配置结构只注册一次；字段 Key 应稳定。默认值来自标签，动态初始值在 `Init()` 中填充。Init 每次启动执行，但已有数据库值不被默认值覆盖。不要在 Init 内使用 Get/MustGet，亦不需要手动调用框架 BaseConfig 初始化。

每次 Get 读取指定 Key；命中缓存不查数据库，未命中按项查询并回填，缓存异常时可回源数据库。配置默认缓存 30 分钟，普通读取不会自动续期。数据库更新与缓存更新不是跨存储原子事务；缓存写入失败会尝试移除旧值并返回错误。

## 缓存与锁

`admin.Cache()` 返回内存或 Redis 的统一接口；显式声明类型可用 `admin.CacheStore` 和 `admin.CacheLock`。

- 写入：`Set(key, value, ttl)`、`SetJSON(key, value, ttl)`；ttl 为 0 永不过期。
- 读取：`Get`、`GetBool`、`GetInt`、`GetJSON(key, &target)`。
- 管理：`Exists(key)`、`Del(key)`，未命中可用 `errors.Is(err, admin.ErrCacheNotFound)` 判断。
- `MustGet / MustGetBool / MustGetInt` 忽略读取错误；**不同于配置项 MustGet 的 panic 语义**。
- `Lock(ctx, key, expiration, extendInterval)` 尝试一次；`WaitForLock(...)` 等待并重试，务必给 ctx 设置超时。
- 获取成功后调用 `Unlock()`；自动续期不能保证网络故障下业务绝不重叠。

```go
store := admin.Cache()
if err := store.Set("product:featured", "42", 30*time.Minute); err != nil {
    return err
}
value, err := store.Get("product:featured")
```

框架按数据库驱动、地址、端口、数据库名的 MD5 建立命名空间，再区分缓存、锁等用途；业务只传 `模块:用途:标识`。`Client()` 返回原始客户端，不自动加命名空间，应优先使用抽象接口。

## 注册计划任务处理函数

在模块 `Initialize()` 中注册处理函数，随后在计划任务页面创建任务实例、填写执行周期和参数：

```go
return admin.Scheduler().RegisterHandler("product.cleanup", &admin.HandlerOption{
    Name: "清理过期商品",
    Params: admin.HandlerParams{
        &admin.HandlerParam{
            Name: "保留天数", Key: "days", Type: admin.IntParams, Required: true,
        },
    },
    Handler: func(ctx context.Context, raw []byte) error {
        var params struct { Days int `json:"days"` }
        if err := json.Unmarshal(raw, &params); err != nil { return err }
        // 校验参数并执行业务；通过 ctx 响应取消。
        return nil
    },
})
```

参数类型为 `StringParams、IntParams、BoolParams`，处理函数仍应验证业务范围。`GetHandlers()` 可读取已注册处理函数。注册处理函数不会自动创建计划任务记录。

默认时区为北京时间；五段 Cron 示例 `0 3 * * *` 表示每天 03:00。框架负责 StartScheduler/StopScheduler，业务不要再次启动调度器。多实例使用数据库任务领取和共享缓存锁协调，业务任务仍需幂等，不能依赖“绝对只执行一次”。

## 全局 Gin 配置

```go
err := admin.ConfigureEngine(func(engine *gin.Engine) error {
    // 例如 engine.Use(cors.New(...))，参数由业务部署环境决定。
    return engine.SetTrustedProxies([]string{"127.0.0.1"})
})
if err != nil { log.Fatal(err) }
// 然后 admin.Run()
```

回调在模块初始化完成后、框架路由注册前执行。它适合全局中间件和代理配置；直接在 engine 注册的接口不自动参与框架鉴权及 API 同步。返回错误会中止启动。

## 部署与升级注意事项

- 单实例可使用内存缓存；多实例必须使用同一数据库、共享 Redis，并保持数据库命名空间参数和 JWT 密钥一致。
- 密码变更、管理员状态及角色关系通过框架接口修改，才能同步维护身份缓存和会话失效；直接写表不会触发这些逻辑。
- Casbin 在各进程维护权限，Redis 用于变更通知，后台重新同步用于恢复；不要把它理解为每次请求都读权限表。
- 生产关闭 dev，初始化后修改默认账号密码；代理部署时设置可信代理。
- 升级使用 `go get github.com/Gary-Yez/go-admin@vX.Y.Z`，随后 `go mod tidy` 并构建。配套前端使用对应发行版本，升级前检查发布说明与数据库变更。

公共 API 以根包 `admin`、`request`、`response` 为入口。`internal/` 是框架实现，不应复制或直接依赖其中的系统模块与服务。
