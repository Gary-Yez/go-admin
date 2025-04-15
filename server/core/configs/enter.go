package configs

var Config *config

type config struct {
	Server server      `mapstructure:"server"`
	Mysql  mysqlConfig `mapstructure:"mysql"`
	Redis  redisConfig `mapstructure:"redis"`
	Jwt    jwtConfig   `mapstructure:"jwt"`
}

func init() {
	Config = &config{}
}
