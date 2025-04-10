package configs

type listen struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
type deploy struct {
	AdminPrefix string  `yaml:"admin_prefix"`
	Listen      *listen `yaml:"listen"`
}
