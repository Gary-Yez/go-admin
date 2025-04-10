package configs

type listen struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type deploy struct {
	AdminPrefix string  `yaml:"adminPrefix"`
	Listen      *listen `yaml:"listen"`
}
