package config

type Config struct {
	Server Server `mapstructure:"server"`
	Mysql  Mysql  `mapstructure:"mysql"`
	Redis  Redis  `mapstructure:"redis"`
	Jwt    Jwt    `mapstructure:"jwt"`
}

func (c *Config) IsDev() bool {
	return c.Server.Dev
}
