package structs

type BucketConfig struct {
	Endpoint  string `yaml:"endpoint"`
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
	UseSSL    bool   `yaml:"useSSL"`
}
