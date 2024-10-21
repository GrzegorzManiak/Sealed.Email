package structs

type DomainServiceConfig struct {
	RetryMax      int `yaml:"retryMax"`
	RetryEvery    int `yaml:"retryEvery"`
	MaxConcurrent int `yaml:"maxConcurrent"`
	BatchSize     int `yaml:"batchSize"`
	BatchTimeout  int `yaml:"batchTimeout"`
}

type DomainConfig struct {
	ProvidedDomains  []string            `yaml:"provided"`
	IdentityProvider string              `yaml:"identityProvider"`
	SpfRecord        string              `yaml:"spfRecord"`
	MxRecords        []string            `yaml:"mxRecords"`
	DkimTemplate     string              `yaml:"dkimTemplate"`
	DmarcTemplate    string              `yaml:"dmarcTemplate"`
	Service          DomainServiceConfig `yaml:"service"`
}
