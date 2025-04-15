package configs

type server struct {
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	AdminPrefix string `mapstructure:"admin_prefix"`
}
