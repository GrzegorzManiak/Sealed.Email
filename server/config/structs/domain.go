package structs

type DomainServiceConfig struct {
	RetryMax      int64 `yaml:"retryMax"`
	RetryInterval int64 `yaml:"retryInterval"`
	MaxConcurrent int   `yaml:"maxConcurrent"`
	BatchSize     int   `yaml:"batchSize"`
	BatchTimeout  int   `yaml:"batchTimeout"`
}

type DomainConfig struct {
	ProvidedDomains   []string `yaml:"provided"`
	IdentityTemplate  string   `yaml:"identityTemplate"`
	SpfRecord         string   `yaml:"spfRecord"`
	MxRecords         []string `yaml:"mxRecords"`
	DkimTemplate      string   `yaml:"dkimTemplate"`
	DmarcTemplate     string   `yaml:"dmarcTemplate"`
	DKIMKeySize       int      `yaml:"DKIMKeySize"`
	ChallengeTemplate string   `yaml:"challengeTemplate"`

	Service DomainServiceConfig `yaml:"service"`
}
