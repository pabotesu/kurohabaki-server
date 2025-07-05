package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestConfigStructure tests the config structure directly
func TestConfigStructure(t *testing.T) {
	// Create a test configuration directly
	cfg := Config{
		ServerSettings: ServerSettings{
			InterfaceName: "test0",
			PrivateKey:    "YFRUWNhFCHmLn3WYxJBBXvJlVeynEJlfzYQsTtZQPWQ=",
			PublicIP:      "203.0.113.1",
			KhIP:          "10.0.0.1",
			Port:          51820,
			KhNwRange:     "10.0.0.0/24",
			Etcd: EtcdConfig{
				ListenClientIP:      "127.0.0.1",
				ListenClientPort:    2379,
				AdvertiseClientIP:   "203.0.113.1",
				AdvertiseClientPort: 2379,
				EndpointIP:          "127.0.0.1",
				EndpointPort:        2379,
			},
		},
	}

	// Verify the config
	if cfg.ServerSettings.InterfaceName != "test0" {
		t.Errorf("Expected InterfaceName to be 'test0', got '%s'", cfg.ServerSettings.InterfaceName)
	}

	if cfg.ServerSettings.KhNwRange != "10.0.0.0/24" {
		t.Errorf("Expected KhNwRange to be '10.0.0.0/24', got '%s'", cfg.ServerSettings.KhNwRange)
	}

	if cfg.ServerSettings.Port != 51820 {
		t.Errorf("Expected Port to be 51820, got %d", cfg.ServerSettings.Port)
	}

	if cfg.ServerSettings.Etcd.AdvertiseClientPort != 2379 {
		t.Errorf("Expected Etcd.AdvertiseClientPort to be 2379, got %d", cfg.ServerSettings.Etcd.AdvertiseClientPort)
	}
}

// TestLoadConfig tests the config loading functionality
func TestLoadConfig(t *testing.T) {
	// テスト用の一時設定ファイルを作成
	tmpDir, err := os.MkdirTemp("", "kurohabaki-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
server_settings:
  interface_name: "test0"
  private_key: "YFRUWNhFCHmLn3WYxJBBXvJlVeynEJlfzYQsTtZQPWQ="
  public_ip: "203.0.113.1"
  kh_ip: "10.0.0.1"
  port: 51820
  kh_nw_range: "10.0.0.0/24"
  etcd:
    listen_client_ip: "127.0.0.1"
    listen_client_port: 2379
    advertise_client_ip: "203.0.113.1"
    advertise_client_port: 2379
    etcd_endpoint_ip: "127.0.0.1"
    etcd_endpoint_port: 2379
`

	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// 設定を読み込み
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// 設定値を検証
	if cfg.ServerSettings.InterfaceName != "test0" {
		t.Errorf("Expected InterfaceName to be 'test0', got '%s'", cfg.ServerSettings.InterfaceName)
	}

	if cfg.ServerSettings.Port != 51820 {
		t.Errorf("Expected Port to be 51820, got %d", cfg.ServerSettings.Port)
	}

	// 存在しない設定ファイルのテスト
	_, err = LoadConfig("/nonexistent/config.yaml")
	if err == nil {
		t.Error("Expected error when loading non-existent config file")
	}
}

// TestInvalidConfig tests error handling for invalid configurations
func TestInvalidConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kurohabaki-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 無効なYAML形式
	invalidConfigPath := filepath.Join(tmpDir, "invalid.yaml")
	invalidContent := `
server_settings:
  interface_name: "test0"
  port: not_a_number  # ここが無効
`

	err = os.WriteFile(invalidConfigPath, []byte(invalidContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid config: %v", err)
	}

	_, err = LoadConfig(invalidConfigPath)
	if err == nil {
		t.Error("Expected error when loading invalid config file")
	}
}
