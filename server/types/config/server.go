package config

type Server struct {
	Dev         bool   `mapstructure:"dev"`
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	AdminPrefix string `mapstructure:"admin_prefix"`
	ApiPrefix   string `mapstructure:"api_prefix"`
}
