package cmd

import "testing"

func TestNodeCmd(t *testing.T) {
	// Test basic command structure
	if got := nodeCmd.Use; got != "node" {
		t.Errorf("nodeCmd.Use = %q, expected %q", got, "node")
	}

	if got := nodeCmd.Short; got != "Manage nodes" {
		t.Errorf("nodeCmd.Short = %q, expected %q", got, "Manage nodes")
	}

	// Test subcommand registration
	foundNodeAddCmd := false
	for _, cmd := range nodeCmd.Commands() {
		if cmd == nodeAddCmd {
			foundNodeAddCmd = true
			break
		}
	}

	if !foundNodeAddCmd {
		t.Error("nodeAddCmd is not registered as a subcommand of nodeCmd")
	}

	// Test that nodeCmd is registered to rootCmd
	foundNodeCmd := false
	for _, cmd := range rootCmd.Commands() {
		if cmd == nodeCmd {
			foundNodeCmd = true
			break
		}
	}

	if !foundNodeCmd {
		t.Error("nodeCmd is not registered as a subcommand of rootCmd")
	}
}
