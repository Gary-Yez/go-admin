package config

type Jwt struct {
	Secret string `mapstructure:"secret"`
}
