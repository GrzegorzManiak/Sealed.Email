package structs

type SmtpQueueConfig struct {
	Name          string `yaml:"name"`
	RetryMax      int64  `yaml:"retryMax"`
	RetryInterval int64  `yaml:"retryInterval"`
	MaxConcurrent int    `yaml:"maxConcurrent"`
	BatchSize     int    `yaml:"batchSize"`
	BatchTimeout  int    `yaml:"batchTimeout"`
}

type SmtpConfig struct {
	OutboundQueue SmtpQueueConfig `yaml:"outboundQueue"`
	InboundQueue  SmtpQueueConfig `yaml:"inboundQueue"`
}
