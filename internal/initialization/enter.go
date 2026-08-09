package initialization

import (
	"errors"
	"github.com/Gary-Yez/go-admin/cache"
	"github.com/Gary-Yez/go-admin/config"
	"github.com/Gary-Yez/go-admin/scheduler"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func InitConfig(flag *pflag.FlagSet, configFile string) (*config.Config, error) {
	cfg := &config.Config{}
	viper.SetEnvPrefix("MYAPP")
	viper.AutomaticEnv()
	_ = viper.BindPFlag("server.host", flag.Lookup("server.host"))
	_ = viper.BindPFlag("server.port", flag.Lookup("server.port"))
	viper.SetConfigFile(configFile)
	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

type Dependencies struct {
	DB        *gorm.DB
	Cache     cache.Cache
	Scheduler scheduler.Scheduler
}

func InitDependencies(cfg *config.Config) (*Dependencies, error) {
	db, err := initGormMysql(cfg)
	if err != nil {
		return nil, errors.New("数据库初始化失败：" + err.Error())
	}
	cacheStore, err := initCache(cfg)
	if err != nil {
		return nil, errors.New("缓存初始化失败：" + err.Error())
	}
	scheduler, err := initTaskManager(cacheStore)
	if err != nil {
		return nil, errors.New("任务管理器初始化失败：" + err.Error())
	}
	return &Dependencies{DB: db, Cache: cacheStore, Scheduler: scheduler}, nil
}
