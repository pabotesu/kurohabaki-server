package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"
)

// Reset the test environment
func resetTestEnv() {
	configPath = ""
	// Don't re-register flags - use already registered ones
}

// Set flags once before test execution
func init() {
	// Register flags for testing
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to the configuration file (optional)")
}

func TestGetConfigPath(t *testing.T) {
	resetTestEnv()

	t.Run("Initial state", func(t *testing.T) {
		if got := GetConfigPath(); got != "" {
			t.Errorf("GetConfigPath() = %q, expected empty string", got)
		}
	})

	t.Run("Set value", func(t *testing.T) {
		expected := "/path/to/config.yaml"
		configPath = expected
		if got := GetConfigPath(); got != expected {
			t.Errorf("GetConfigPath() = %q, expected %q", got, expected)
		}
	})
}

func TestConfigFlag(t *testing.T) {
	resetTestEnv()

	t.Run("Parse config flag", func(t *testing.T) {
		testPath := "/test/config.yaml"

		// Parse the flag
		cmd := rootCmd
		// Don't re-register flags - they're already registered

		err := cmd.ParseFlags([]string{"--config", testPath})
		if err != nil {
			t.Fatalf("Failed to parse flags: %v", err)
		}

		// Check if flag value is set correctly
		if got := configPath; got != testPath {
			t.Errorf("configPath = %q, expected %q", got, testPath)
		}
	})
}

func TestPersistentPreRun(t *testing.T) {
	resetTestEnv()

	t.Run("PersistentPreRun output", func(t *testing.T) {
		testPath := "/test/config.yaml"
		configPath = testPath

		// Capture stdout
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Execute PersistentPreRun
		rootCmd.PersistentPreRun(rootCmd, []string{})

		w.Close()
		os.Stdout = oldStdout

		// Read captured output
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r)
		if err != nil {
			t.Fatalf("Failed to read captured output: %v", err)
		}

		// Verify output
		expected := "Using config file: " + testPath + "\n"
		if got := buf.String(); got != expected {
			t.Errorf("PersistentPreRun output = %q, expected %q", got, expected)
		}
	})
}

// Test basic settings of rootCmd
func TestRootCmdBasics(t *testing.T) {
	// Test command name
	if got := rootCmd.Use; got != "kurohabaki-server" {
		t.Errorf("rootCmd.Use = %q, expected %q", got, "kurohabaki-server")
	}

	// Test short description
	expectedShort := "kurohabaki-server is a WireGuard-based P2P coordination server"
	if got := rootCmd.Short; got != expectedShort {
		t.Errorf("rootCmd.Short = %q, expected %q", got, expectedShort)
	}

	// Test if config flag is registered
	flag := rootCmd.PersistentFlags().Lookup("config")
	if flag == nil {
		t.Error("'config' flag is not registered")
		return
	}

	if flag.Usage != "Path to the configuration file (optional)" {
		t.Errorf("Config flag usage differs: %q", flag.Usage)
	}
}
