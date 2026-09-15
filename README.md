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

## 文件存储抽象

公开的 storage 包提供统一 storage.New(config) 创建入口、包含普通上传、读取、删除和分片操作的 Store 接口，以及可选 URLProvider。已提供本地、腾讯云 COS、阿里云 OSS、S3 四种存储适配器，R2 和 MinIO 统一通过 S3 配置接入，可同时使用多个实例；支持逐片上传、合并和取消，内置文件管理模块负责上传 HTTP 接口及会话持久化。

storage.New 只支持 Local、Tencent、Aliyun、S3 四个内置引擎。S3 兼容服务通过 Endpoint、Region、Bucket、凭据和 PathStyle 接入，不需要注册构造函数。未知引擎直接返回错误。

所有内置适配器都实现完整的 Store 接口，可以直接调用 BeginMultipart、UploadPart、CompleteMultipart、AbortMultipart，无需类型断言。供应商 SDK 原始错误继续返回，创建失败时实例为 nil；本地适配器用完后通过 io.Closer 关闭。

storage 包只处理配置和 SDK 操作；存储账号、缓存、权限、文件记录及上传会话由系统模块负责。New 每次创建新实例，不缓存适配器；云存储继续复用 HTTP 连接池。S3 的 CredentialsProvider 用于 SDK 凭据获取，独立于自定义引擎注册，仍保留。

## 存储账号与业务调用

“文件管理”下包含“文件列表”和“存储管理”。存储管理可为同一引擎新增多个账号，分别设置名称、桶、服务地址和凭据，并选择默认账号。新安装先创建账号；不自动创建空的云存储配置。

框架完成初始化后，业务可以创建默认存储或指定账号：

```go
store, err := admin.Storage()       // 已启用的默认存储
store, err := admin.Storage(3)      // 存储管理中的账号 ID
config, err := admin.StorageConfig(3) // 只读取配置，不创建适配器
```

每次 Storage 调用都创建新适配器，配置按账号 ID 整体缓存 30 分钟，复用现有内存/Redis 缓存抽象。修改和删除账号与缓存未命中回填共用锁，更新前失效旧缓存，避免旧值重新回填。默认账号 ID 从数据库解析；S3 兼容服务通过同一个 S3 引擎配置接入。

文件上传先在短数据库事务内核对账号并建立引用，防止账号同时被删除或更换位置，这一步读取数据库当前配置；文件传输在事务结束后进行。后续分片、下载和删除按保存的账号 ID 读取缓存配置。业务直接使用 admin.Storage 上传的对象不自动成为文件管理记录，需要自行管理其关联和生命周期。

已有文件或未完成会话引用的账号不允许更换引擎、桶、区域、地址、路径模式或本地目录，也不允许删除。可更新凭据、名称和下载域名；停用只阻止新上传，不中断已有文件访问及续传。最多一个默认账号，默认账号必须启用；没有默认账号时上传请求需明确提供 storage_id。

后台存储账号只配置长期访问密钥。编辑接口不返回 Secret 原文：留空保留原值，填写新值替换。后台及存储适配器均不提供临时凭据令牌配置。配置保存在数据库和缓存，仍需保护数据库及备份。文件记录不再重复保存凭据。

文件只关联 storage_id，不保存账号配置快照。每次操作使用该账号的最新配置（通过缓存读取，修改账号时失效缓存），没有旧快照迁移或回退逻辑。分片会话只保存上传 ID、对象 Key、大小等续传信息，不包含账号凭据。

## 文件管理与上传接口

启动后自动注册“文件管理 → 文件列表 / 存储管理”，创建文件记录表和分片回执表，API 同步仍遵循现有角色授权。先在存储管理中添加存储账号，再授予角色文件管理菜单及所需接口权限。上传接口使用现有 JWT / API Token 鉴权。

以下路径相对于 API 前缀：

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| POST | /sys_file/list | ReqList 分页、筛选和排序，同时返回上传大小策略 |
| GET | /sys_file/options | 返回上传策略和存储账号选项，不查询文件列表；需要该接口权限 |
| POST | /sys_file/upload | multipart/form-data，file 字段及可选 storage_id，上限由“存储设置”的分片上传阈值决定（默认 100 MiB） |
| POST | /sys_file/begin | JSON：name、size、last_modified、可选 storage_id，创建大文件上传会话 |
| GET | /sys_file/session?id=1 | 查询本人会话和已收到的分片编号 |
| POST | /sys_file/part?id=1&number=1 | application/octet-stream 原始分片，必须携带准确 Content-Length |
| POST | /sys_file/complete | JSON：id，服务端读取已保存回执并合并，支持完成请求重试 |
| POST | /sys_file/abort | JSON：id，取消本人未完成会话并删除分片 |
| POST | /sys_file/delete | JSON：ids，删除记录及原存储中的文件/分片 |
| POST | /sys_file/cleanup | 按各会话的到期时间清理未完成上传 |
| POST | /sys_file/link | JSON：id、preview（可选），取得公开直连地址、按配置有效期生成的签名地址或本地下载凭证 |
| GET | /sys_file/content/:ticket | 使用下载凭证流式获取文件，无需附加 JWT |

配置管理的“存储设置”提供 storage.multipart_threshold_mib（分片阈值，默认 100，范围 100～5120 MiB）、storage.part_size_mib（分片大小，默认 20，范围 20～1024 MiB）和 storage.session_days（上传会话有效期，默认 1 天，范围 1～365 天）。通过已有业务配置缓存读取，修改后影响新上传。已有分片会话继续使用 upload_json 中的 PartSize，前后端都按会话大小续传。会话创建时将有效期转换为 expires_at，后续修改配置不改变已有会话的到期时间；成功上传的文件不受该有效期影响。

普通上传成功返回文件记录。分片流程为 begin → 按编号上传 part → complete；默认超过 100 MiB 自动分片，每片 20 MiB，最后一片可较小，最大 10000 片。begin 返回 file.id 和 part_size；后续请求只提交文件记录 ID，不允许客户端指定供应商 UploadId、对象路径或存储凭据。last_modified 是原文件修改时间的毫秒时间戳；页面重新选文件时核对名称、大小和修改时间，业务不应把这种检查当成完整内容哈希校验。

文件模块注册 sys_clear_file_uploads（文件管理-清理过期上传）处理函数，不需要任务参数。可在计划任务页面选择该处理函数并设置执行周期；注册处理函数不会自动创建或启用任务。清理依据各会话已保存的 expires_at，删除过期未完成上传的存储对象、分片及数据库记录；已完成文件不受影响。存储清理失败会返回错误并保留对应记录供重试。

会话与回执保存在数据库，刷新页面或重启服务后可以继续。重传同编号分片会覆盖该分片；complete 重试先查询最终对象，处理存储已合并、数据库提交或响应失败的情况。数据库行锁串行协调同一文件的分片、合并、删除，不依赖单机内存锁；存储请求期间会占用该事务连接。

每个文件保存 storage_id，续传和旧文件操作始终使用原账号，不受默认账号变化影响。更新同一账号的访问凭据会作用于后续文件操作。


预览和下载优先使用账号配置的 PublicURL，返回公开直连地址，expires_in=0 表示应用不设置过期时间。未配置 PublicURL 的云存储返回 按配置有效期生成的签名地址，文件流量直接经过存储服务，不经过 Go；签名下载携带原文件名。仅本地存储未配置 PublicURL 时使用 Go 临时凭证和流式返回，本地文件支持 Range。link 返回的 HTTP(S) 绝对地址直接使用，本地相对路径由 API 基础地址解析。私有签名和临时凭证在有效期内可重复使用，不应作为永久业务地址保存。

全局上传扩展名由存储设置 storage.allowed_extensions 控制，类型为字符串列表，默认空列表不限制。使用现有字符串列表控件，逐项输入 jpg、png、pdf、zip 并按回车添加，不区分大小写，允许带前导点并自动去重；仅支持单个扩展名，不接受通配符、MIME 或复合后缀。前端筛选并提示，后端在普通上传、创建分片、继续分片及合并时强制检查，未完成会话遵循最新规则，已完成文件不受影响。此设置仅检查文件名扩展名，不替代业务的文件内容校验。

文件列表支持 JPEG、PNG、GIF、WebP、AVIF 图片缩略图及大图，浏览器可播放的 MP4、WebM、OGG 视频和 MP3、M4A、AAC、WAV、OGG、FLAC、WebM 音频，PDF，以及纯文本、CSV、JSON、XML、HTML 和源代码的文本预览。具体准入根据服务端检测的 MIME 判断；媒体能否播放还取决于浏览器编码支持，PDF 依赖浏览器内置阅读器。文本按 UTF-8 展示前 256 KiB，需要存储允许站点跨域读取；HTML、脚本仅显示文本、不执行。Office 文档和压缩包显示类型图标并支持下载，不接入第三方预览。公开域名必须已经由存储服务或 CDN 提供对应对象访问；填写地址不会自动公开桶或挂载本地目录。直连时的 Content-Type、Content-Disposition 和缓存策略由存储服务/CDN 控制，公开地址不保证强制下载或保留原文件名；请按需要配置响应头。云端签名使用的服务地址必须能被浏览器访问。

单机可以使用本地存储；多实例使用同一数据库、共享 Redis 和云存储，各节点需要能访问对应账号的存储地址。本地磁盘不支持跨节点访问，Docker 部署须持久化上传目录及其同级的 .multipart 目录。反向代理请求体上限需覆盖普通上传阈值加 1 MiB 表单开销，并覆盖实际分片大小（默认至少 101 MiB），并为上传、合并设置合理超时。

删除文件会让业务中的引用失效；批量删除可能部分成功，失败记录会保留供重试。删除确认支持“仅删除记录”，对应 records_only=true：只删除数据库文件记录和分片回执，不读取账号配置、不调用存储服务；远端对象或分片需要自行清理。默认仍同时删除存储对象，失败时可以在同一确认框勾选仅删除记录后重试。云存储不参与数据库事务：进程在供应商创建会话成功、UploadId 尚未落库时退出，仍可能遗留供应商分片，建议在桶上配置未完成上传生命周期清理。系统的“清理过期上传”负责已记录会话，不声称能找回未知 UploadId。

存储管理 API（均参与角色授权）：

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| POST | /sys_storage/list | 分页、筛选、排序 |
| GET | /sys_storage/get?id=1 | 编辑详情，凭据脱敏 |
| POST | /sys_storage/save | 新建或编辑账号；id=0 表示新建 |
| POST | /sys_storage/default | 设置默认账号，JSON：id；账号必须已启用 |
| POST | /sys_storage/enabled | 切换启用状态，JSON：id、enabled；停用时清除默认标记 |
| POST | /sys_storage/delete | 删除无文件引用的账号，JSON：ids |

文件列表响应中的 storages 只包含账号 ID、名称、引擎、启用和默认状态，上传者不需要额外取得存储管理权限。

访问链接有效期由存储设置 storage.link_expire_seconds 控制，默认 3600 秒（1 小时），范围 60～86400 秒。私有签名和本地访问凭证使用该有效期，接口 expires_in 返回秒数；修改仅影响新生成的链接，公开地址不受影响。

本地存储未配置公开地址时，同一文件的预览与下载共用一个当前有效的随机凭证。地址为 /sys_file/content/文件ID.凭证，添加 ?download=1 强制下载，不带参数时可预览类型内联展示，其余类型仍下载。缓存只按文件 ID 保存，命中不续期，过期后通过现有缓存锁重新生成；接口 expires_in 返回剩余秒数。访问会核对凭证与绝对到期时间。云存储签名及公开地址行为不变。内存缓存重启或缓存被清理后需要重新获取链接；升级前旧格式的本地临时链接需刷新页面重新获取。
