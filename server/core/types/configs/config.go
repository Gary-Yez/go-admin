package configs

type Config struct {
	Server *server      `yaml:"server"`
	Mysql  *mysqlConfig `yaml:"mysql"`
	Redis  *redisConfig `yaml:"redis"`
	Jwt    *JwtConfig   `yaml:"jwt"`
}
