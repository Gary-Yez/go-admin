package configs

type jwtConfig struct {
	Secret string `mapstructure:"secret"`
}
