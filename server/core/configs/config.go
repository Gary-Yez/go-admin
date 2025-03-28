package configs

type Config struct {
	Mysql *mysqlConfig `yaml:"mysql"`
	Redis *redisConfig `yaml:"redis"`
	Jwt   *jwtConfig   `yaml:"jwt"`
}
