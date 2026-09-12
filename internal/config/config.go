package config

type JWT struct {
	Secret string `mapstructure:"secret"`
}

type Config struct {
	JWT      JWT      `mapstructure:"jwt"`
	Server   Server   `mapstructure:"server"`
	Database Database `mapstructure:"database"`
	Redis    Redis    `mapstructure:"redis"`
}

func (c *Config) IsDev() bool {
	return c.Server.Dev
}
