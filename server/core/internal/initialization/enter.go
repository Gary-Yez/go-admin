package initialization

import (
	"errors"
	"gitee.com/mxcker/go-admin/server/configs"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func InitConfig(flag *pflag.FlagSet) error {
	viper.SetEnvPrefix("MYAPP")
	viper.AutomaticEnv()
	_ = viper.BindPFlag("server.host", flag.Lookup("server.host"))
	_ = viper.BindPFlag("server.port", flag.Lookup("server.port"))
	_ = viper.BindPFlag("config", flag.Lookup("config"))
	viper.SetConfigFile(viper.GetString("config"))
	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&configs.Config); err != nil {
		return err
	}
	return nil
}

func InitGlobal() error {
	err := initGormMysql()
	if err != nil {
		return errors.New("数据库初始化失败：" + err.Error())
	}
	err = initCache()
	if err != nil {
		return errors.New("缓存初始化失败：" + err.Error())
	}
	err = initTaskManager()
	if err != nil {
		return errors.New("任务管理器初始化失败：" + err.Error())
	}
	return nil
}
