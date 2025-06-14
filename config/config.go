package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// EtcdConfig holds the etcd-related configuration.
type EtcdConfig struct {
	ListenClientIP      string `mapstructure:"listen-client-ip"`
	ListenClientPort    int    `mapstructure:"listen-client-port"`
	AdvertiseClientIP   string `mapstructure:"advertise-client-ip"`
	AdvertiseClientPort int    `mapstructure:"advertise-client-port"`
	EndpointIP          string `mapstructure:"etcd-endpoint-ip"`
	EndpointPort        int    `mapstructure:"etcd-endpoint-port"`
}

// ServerSettings contains the main server configuration.
type ServerSettings struct {
	InterfaceName string     `mapstructure:"interface_name"`
	PrivateKey    string     `mapstructure:"private_key"`
	PublicIP      string     `mapstructure:"public_ip"`
	KhIP          string     `mapstructure:"kh_ip"`
	Port          int        `mapstructure:"port"`
	KhNwRange     string     `mapstructure:"kh_nw_range"`
	Etcd          EtcdConfig `mapstructure:"etcd"`
}

// Config is the root of the loaded configuration.
type Config struct {
	ServerSettings ServerSettings `mapstructure:"server_settings"`
}

// LoadConfig reads the configuration file and returns a Config struct.
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()

	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("/etc/kurohabaki/")
		home, err := os.UserHomeDir()
		if err == nil {
			v.AddConfigPath(filepath.Join(home, ".config", "kurohabaki"))
		}
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
