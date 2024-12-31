package structs

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Id   string `yaml:"id"`
}
