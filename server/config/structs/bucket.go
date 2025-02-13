package structs

type BucketConfig struct {
	Key    string `yaml:"key"`
	Api    string `yaml:"api"`
	Secret string `yaml:"secret"`
	UseSsl bool   `yaml:"useSSL"`
}
