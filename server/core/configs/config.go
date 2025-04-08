package configs

type Config struct {
	Server *Server      `yaml:"server"`
	Mysql  *mysqlConfig `yaml:"mysql"`
	Redis  *redisConfig `yaml:"redis"`
	Jwt    *jwtConfig   `yaml:"jwt"`
}
