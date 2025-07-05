package clientconfig

import (
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestClientYAMLConfig(t *testing.T) {
	// Create a config instance with the actual struct definition
	cfg := ClientYAMLConfig{}

	// Set individual fields
	cfg.Interface.PrivateKey = "privatekey123"
	cfg.Interface.Address = "10.0.0.3/32"
	cfg.Interface.DNS = "10.0.0.1"
	cfg.Interface.Routes = []string{"10.0.0.0/24"}

	cfg.Peer.PublicKey = "publickey123"
	cfg.Peer.Endpoint = "203.0.113.1:51820"
	cfg.Peer.AllowedIPs = "10.0.0.1/32"
	cfg.Peer.PersistentKeepalive = 5

	cfg.Etcd.Endpoint = "203.0.113.1:2379"

	// Serialize to YAML
	yamlData, err := yaml.Marshal(&cfg)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	// Deserialize back from YAML
	var unmarshalled ClientYAMLConfig
	err = yaml.Unmarshal(yamlData, &unmarshalled)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Verify the config matches the original
	if !reflect.DeepEqual(cfg, unmarshalled) {
		t.Errorf("YAML roundtrip failed. Original: %+v, Unmarshalled: %+v", cfg, unmarshalled)
	}
}

func TestClientYAMLConfigBasic(t *testing.T) {
	// Test with minimal configuration
	cfg := ClientYAMLConfig{}

	// Set individual fields
	cfg.Interface.PrivateKey = "privatekey123"
	cfg.Interface.Address = "10.0.0.3/32"

	// Validate field values
	if cfg.Interface.PrivateKey != "privatekey123" {
		t.Errorf("Expected PrivateKey to be 'privatekey123', got '%s'", cfg.Interface.PrivateKey)
	}

	// Validate YAML tags
	yamlData, err := yaml.Marshal(&cfg)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	yamlStr := string(yamlData)
	t.Logf("Generated YAML:\n%s", yamlStr)

	// Check if the YAML contains expected strings
	expectedStrings := []string{
		"private_key: privatekey123",
		"address: 10.0.0.3/32",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(yamlStr, expected) {
			t.Errorf("Expected YAML to contain '%s', but it didn't", expected)
		}
	}
}

func TestYAMLTags(t *testing.T) {
	// Verify YAML tags function correctly
	cfg := ClientYAMLConfig{}
	cfg.Interface.PrivateKey = "test-key"
	cfg.Peer.PersistentKeepalive = 25

	// Serialize to YAML
	yamlData, err := yaml.Marshal(&cfg)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	yamlStr := string(yamlData)

	// Check for expected YAML key names
	expectedKeys := []string{
		"interface:",
		"private_key:",
		"peer:",
		"public_key:",
		"persistent_keepalive:",
		"etcd:",
	}

	for _, key := range expectedKeys {
		if !strings.Contains(yamlStr, key) {
			t.Errorf("Expected YAML to contain key '%s', but it didn't", key)
		}
	}

	// Also check values
	if !strings.Contains(yamlStr, "test-key") {
		t.Error("Expected private_key value 'test-key' not found in YAML")
	}

	if !strings.Contains(yamlStr, "persistent_keepalive: 25") {
		t.Error("Expected persistent_keepalive value '25' not found in YAML")
	}
}
