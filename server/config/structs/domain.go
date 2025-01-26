package structs

type DomainServiceConfig struct {
	MaxRetry      int64  `yaml:"maxRetry"`
	RetryInterval int64  `yaml:"retryInterval"`
	MaxConcurrent int    `yaml:"maxConcurrent"`
	BatchSize     int    `yaml:"batchSize"`
	BatchTimeout  int    `yaml:"batchTimeout"`
	DNS           string `yaml:"dns"`
	VerifyAll     bool   `yaml:"verifyAll"`
}

type DomainConfig struct {
	ProvidedDomains   []string `yaml:"provided"`
	IdentityTemplate  string   `yaml:"identityTemplate"`
	SpfRecordTemplate string   `yaml:"spfRecordTemplate"`
	MxRecords         []string `yaml:"mxRecords"`
	DkimTemplate      string   `yaml:"dkimTemplate"`
	DmarcTemplate     string   `yaml:"dmarcTemplate"`
	DKIMKeySize       int      `yaml:"DKIMKeySize"`
	ChallengeTemplate string   `yaml:"challengeTemplate"`
	ChallengePrefix   string   `yaml:"challengePrefix"`

	Service DomainServiceConfig `yaml:"service"`
}
