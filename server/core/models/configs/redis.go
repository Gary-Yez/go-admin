package configs

import "fmt"

type redisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DB       string `yaml:"db"`
}

func (config *redisConfig) ToString() string {
	dsn := fmt.Sprintf("redis://%s:%s@%s:%s/%s",
		config.Username, config.Password, config.Host, config.Port, config.DB)
	return dsn
}
