package configs

type server struct {
	Dev         bool   `mapstructure:"dev"`
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	AdminPrefix string `mapstructure:"admin_prefix"`
}
