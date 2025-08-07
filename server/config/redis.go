package config

import (
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func (config *Redis) ToString() string {
	dsn := fmt.Sprintf("redis://%s:%s@%s:%s/%s",
		config.Username, config.Password, config.Host, config.Port, config.DB)
	return dsn
}

func (config *Redis) IsNotEmpty() bool {
	return config.Host != "" && config.Port != ""
}

func (config *Redis) Address() string {
	return config.Host + ":" + config.Port
}

func (config *Redis) Option() *redis.Options {
	return &redis.Options{
		Addr:     config.Address(),
		Username: config.Username,
		Password: config.Password, // 没有密码默认留空
		DB:       config.DB,       // 默认使用0号数据库
	}
}
