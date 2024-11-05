package structs

type ServiceConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Prefix   string `yaml:"prefix"`
	TTL      int64  `yaml:"ttl"`
}

type ConnectionPoolConfig struct {
	RefreshInterval int `yaml:"refreshInterval"`
}

type EtcdConfig struct {
	Endpoints      []string             `yaml:"endpoints"`
	ConnectionPool ConnectionPoolConfig `yaml:"connectionPool"`

	Domain       ServiceConfig `yaml:"domain"`
	Notification ServiceConfig `yaml:"notification"`
	SMTP         ServiceConfig `yaml:"smtp"`
	API          ServiceConfig `yaml:"api"`
}
