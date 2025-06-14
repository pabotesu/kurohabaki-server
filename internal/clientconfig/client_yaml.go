package clientconfig

// ClientYAMLConfig defines the structure for the client YAML configuration.
type ClientYAMLConfig struct {
	Interface struct {
		PrivateKey string   `yaml:"private_key"`
		Address    string   `yaml:"address"`
		DNS        string   `yaml:"dns"`
		Routes     []string `yaml:"routes"`
	} `yaml:"interface"`
	Peer struct {
		PublicKey           string `yaml:"public_key"`
		Endpoint            string `yaml:"endpoint"`
		AllowedIPs          string `yaml:"allowed_ips"`
		PersistentKeepalive int    `yaml:"persistent_keepalive"`
	} `yaml:"peer"`
	Etcd struct {
		Endpoint string `yaml:"endpoint"`
	} `yaml:"etcd"`
}
