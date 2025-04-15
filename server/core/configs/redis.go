package configs

import "fmt"

type redisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func (config *redisConfig) ToString() string {
	dsn := fmt.Sprintf("redis://%s:%s@%s:%s/%s",
		config.Username, config.Password, config.Host, config.Port, config.DB)
	return dsn
}
