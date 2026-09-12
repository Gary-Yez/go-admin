package initialization

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Gary-Yez/go-admin/internal/config"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const defaultConfig = `# 修改数据库等配置后，请重新启动程序。
server:
  node_name: "" # 节点显示名称，留空使用主机名；实例 ID 每次启动自动生成
  dev: false # 开发环境需要代码生成工具时改为 true
  host: "0.0.0.0"
  port: "8080"
  admin_prefix: "/admin"
  api_prefix: "/api"
jwt:
  secret: "%s" # 自动随机生成；多实例须使用相同密钥，修改后重启生效
mysql:
  host: "127.0.0.1"
  port: "3306"
  username: ""
  password: ""
  database: ""
redis:
  host: "" # 留空时使用内存缓存
  port: "6379"
  username: ""
  password: ""
  db: 0
`

func InitConfig(flags *pflag.FlagSet, configFile string) (*config.Config, error) {
	if strings.TrimSpace(configFile) == "" {
		return nil, errors.New("配置文件路径不能为空")
	}
	path, err := filepath.Abs(configFile)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件路径失败: %w", err)
	}
	if err := ensureConfigFile(path); err != nil {
		return nil, err
	}
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	v.SetEnvPrefix("MYAPP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	if err := v.BindEnv("server.node_name"); err != nil {
		return nil, err
	}
	if err := v.BindEnv("jwt.secret"); err != nil {
		return nil, err
	}
	v.SetDefault("server.dev", false)
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", "8080")
	for _, key := range []string{"server.host", "server.port"} {
		if flags != nil && flags.Lookup(key) != nil {
			if err := v.BindPFlag(key, flags.Lookup(key)); err != nil {
				return nil, err
			}
		}
	}
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件 %s 失败: %w", path, err)
	}
	cfg := new(config.Config)
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件 %s 失败: %w", path, err)
	}
	if len(strings.TrimSpace(cfg.JWT.Secret)) < 32 {
		return nil, errors.New("配置 jwt.secret 至少需要 32 字节，请填写随机签名密钥后重新启动")
	}
	return cfg, nil
}

func ensureConfigFile(path string) error {
	info, err := os.Stat(path)
	if err == nil {
		if !info.Mode().IsRegular() {
			return fmt.Errorf("配置路径不是普通文件: %s", path)
		}
		return nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("检查配置文件 %s 失败: %w", path, err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}
	// 只创建新文件，避免覆盖已有配置。
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return fmt.Errorf("创建配置文件 %s 失败: %w", path, err)
	}
	_, writeErr := file.WriteString(fmt.Sprintf(defaultConfig, rand.Text()+rand.Text()))
	closeErr := file.Close()
	if err := errors.Join(writeErr, closeErr); err != nil {
		return fmt.Errorf("写入默认配置 %s 失败: %w", path, err)
	}
	return fmt.Errorf("已创建默认配置文件 %s，请修改配置后重新启动程序", path)
}
