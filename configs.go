package admin

import (
	"fmt"
	"reflect"

	"github.com/Gary-Yez/go-admin/internal/bizconfig"
)

// BaseConfig 是框架提供的基础业务配置，环境连接信息仍在 config.yaml。
type BaseConfig struct {
	SiteName                     ConfigItem[string]   `config:"site.name" label:"网站名称" default:"Go Admin" group:"站点设置" description:"登录页及后台 Logo 旁显示的系统名称"`
	SiteTitle                    ConfigItem[string]   `config:"site.title" label:"网站标题" default:"Go Admin 管理后台" group:"站点设置" description:"浏览器标签页显示的标题"`
	SiteLogo                     ConfigItem[string]   `config:"site.logo" label:"网站 Logo" default:"/logo.png" group:"站点设置" description:"支持站内资源路径或完整图片地址"`
	SiteFavicon                  ConfigItem[string]   `config:"site.favicon" label:"浏览器图标" default:"/logo.png" group:"站点设置" description:"浏览器标签页图标，支持站内资源路径或完整图片地址"`
	SiteCopyright                ConfigItem[string]   `config:"site.copyright" label:"版权信息" default:"Copyright © 2026 GaryYez. All rights reserved." group:"站点设置" description:"用于页面底部展示的版权文案，留空表示不展示"`
	JwtExpireMinutes             ConfigItem[int]      `config:"jwt.expire_minutes" label:"登录有效期（分钟）" default:"10080" group:"登录设置" description:"仅影响新签发的登录令牌，默认 7 天"`
	LoginMaxFailures             ConfigItem[int]      `config:"login.max_failures" label:"登录失败次数上限" default:"5" group:"登录设置" description:"同一账号和来源 IP 在统计窗口内达到此次数后临时锁定"`
	LoginFailureWindowMinutes    ConfigItem[int]      `config:"login.failure_window_minutes" label:"失败统计窗口（分钟）" default:"15" group:"登录设置" description:"从首次失败开始统计，窗口过期后重新计数"`
	LoginLockMinutes             ConfigItem[int]      `config:"login.lock_minutes" label:"登录锁定时长（分钟）" default:"15" group:"登录设置" description:"达到失败上限后禁止对应账号和来源 IP 登录的时长"`
	PasswordMinLength            ConfigItem[int]      `config:"login.password_min_length" label:"密码最小长度" default:"8" group:"登录设置" description:"新建账号、修改及重置密码时校验字符数，不影响已有密码登录；密码最多 72 字节"`
	StorageMultipartThresholdMiB ConfigItem[int]      `config:"storage.multipart_threshold_mib" label:"分片上传阈值（MiB）" default:"100" group:"存储设置" description:"超过此大小使用分片上传，范围 100～5120 MiB；仅影响新上传"`
	StoragePartSizeMiB           ConfigItem[int]      `config:"storage.part_size_mib" label:"分片大小（MiB）" default:"20" group:"存储设置" description:"每个分片的大小，范围 20～1024 MiB；最后一片可较小，已有会话保持原大小"`
	StorageSessionDays           ConfigItem[int]      `config:"storage.session_days" label:"上传会话有效期（天）" default:"1" group:"存储设置" description:"未完成上传会话的有效期，范围 1～365 天；仅影响新建会话，完成的文件不受影响"`
	StorageLinkExpireSeconds     ConfigItem[int]      `config:"storage.link_expire_seconds" label:"访问链接有效期（秒）" default:"3600" group:"存储设置" description:"预览和下载链接的有效期，范围 60～86400 秒；仅影响新生成的私有签名和本地访问凭证，公开地址不受影响"`
	StorageAllowedExtensions     ConfigItem[[]string] `config:"storage.allowed_extensions" label:"允许上传的文件类型" default:"[]" group:"存储设置" description:"输入扩展名后按回车添加，例如 jpg、png、pdf、zip；列表为空不限制，不区分大小写；未完成上传按最新规则校验，仅检查扩展名"`
}

type ConfigItem[T bizconfig.Value] struct {
	bizconfig.Item[T]
}

func (item *ConfigItem[T]) Get() (T, error) { return item.Item.Get() }

func (item *ConfigItem[T]) MustGet() T { return item.Item.MustGet() }

// ConfigDefault 用于私有配置结构的 Init 方法，只影响缺失配置的初始值。
func ConfigDefault[T bizconfig.Value](value T) ConfigItem[T] {
	return ConfigItem[T]{Item: bizconfig.InitialValue(value)}
}

// 配置结构只在启动前注册，服务运行期间保持不变。
var configSchema *bizconfig.Schema
var configReady bool

// RegisterConfig 由 settings 包的 init 注册，只绑定 Key，不访问数据库和缓存。
func RegisterConfig[T any](target *T) error {
	if configReady || registry.initialized {
		return fmt.Errorf("配置结构必须在启动服务前注册")
	}
	schema, err := bizconfig.ParseSchema(reflect.TypeFor[T]())
	if err != nil {
		return err
	}
	if configSchema != nil {
		return fmt.Errorf("配置结构已经注册")
	}
	if target == nil {
		return fmt.Errorf("配置对象不能为空")
	}
	if err := schema.Bind(reflect.ValueOf(target).Elem()); err != nil {
		return err
	}
	configSchema = schema
	return nil
}
