package structs

type SmtpQueueConfig struct {
	Name          string `yaml:"name"`
	RetryMax      int64  `yaml:"retryMax"`
	RetryInterval int64  `yaml:"retryInterval"`
	MaxConcurrent int    `yaml:"maxConcurrent"`
	BatchSize     int    `yaml:"batchSize"`
	BatchTimeout  int    `yaml:"batchTimeout"`
}

type SmtpPorts struct {
	Tls      int `yaml:"tls"`
	Plain    int `yaml:"plain"`
	StartTls int `yaml:"startTls"`
}

type SmtpConfig struct {
	OutboundQueue SmtpQueueConfig     `yaml:"outboundQueue"`
	InboundQueue  SmtpQueueConfig     `yaml:"inboundQueue"`
	Certificates  ServiceCertificates `yaml:"certificates"`
	Ports         SmtpPorts           `yaml:"ports"`

	Address           string `yaml:"address"`
	Domain            string `yaml:"domain"`
	WriteTimeout      int    `yaml:"writeTimeout"`
	ReadTimeout       int    `yaml:"readTimeout"`
	MaxMessageBytes   int64  `yaml:"maxMessageBytes"`
	MaxRecipients     int    `yaml:"maxRecipients"`
	AllowInsecureAuth bool   `yaml:"allowInsecureAuth"`
	MaxLineLength     int    `yaml:"maxLineLength"`
	SpfHardFail       bool   `yaml:"spfHardFail"`
}
