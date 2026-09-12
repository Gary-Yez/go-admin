package config

type Server struct {
	NodeName    string `mapstructure:"node_name"`
	Dev         bool   `mapstructure:"dev"`
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	AdminPrefix string `mapstructure:"admin_prefix"`
	ApiPrefix   string `mapstructure:"api_prefix"`
}
