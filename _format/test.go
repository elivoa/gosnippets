type DealTapeConfig struct {
	Endpoint      string `yaml:"end-point"`
	Bid           string `yaml:"bid"`
	User          string `yaml:"user"`
	Pass          string `yaml:"pass"`
	Token         string `yaml:"token"`
	AccessKey     string `yaml:"access-key" json:"access-key"`
	Secret        string `yaml:"secret" json:"secret"`
	Authorization string `yaml:"authorization"`
	StaffId       string `yaml:"staff-id"`
}
