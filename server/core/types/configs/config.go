package configs

type Config struct {
	Deploy *deploy      `yaml:"deploy"`
	Mysql  *mysqlConfig `yaml:"mysql"`
	Redis  *redisConfig `yaml:"redis"`
	Jwt    *JwtConfig   `yaml:"jwt"`
}
