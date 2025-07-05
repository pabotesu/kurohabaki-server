package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

// Only test the command structure and argument validation
func TestNodeAddCmd(t *testing.T) {
	// Test command structure
	if got := nodeAddCmd.Use; got != "add <public-key>" {
		t.Errorf("nodeAddCmd.Use = %q, expected %q", got, "add <public-key>")
	}

	if got := nodeAddCmd.Short; got != "Register a new node by its WireGuard public key" {
		t.Errorf("nodeAddCmd.Short = %q, expected %q", got, "Register a new node by its WireGuard public key")
	}

	// Test Args validation
	if nodeAddCmd.Args == nil {
		t.Fatal("nodeAddCmd.Args is nil, expected cobra.ExactArgs(1)")
	}

	// Test that it requires exactly one argument
	testCmd := &cobra.Command{}
	err := nodeAddCmd.Args(testCmd, []string{})
	if err == nil {
		t.Error("nodeAddCmd.Args did not return an error for empty arguments")
	}

	err = nodeAddCmd.Args(testCmd, []string{"arg1", "arg2"})
	if err == nil {
		t.Error("nodeAddCmd.Args did not return an error for multiple arguments")
	}

	err = nodeAddCmd.Args(testCmd, []string{"valid-pubkey"})
	if err != nil {
		t.Errorf("nodeAddCmd.Args returned an error for valid arguments: %v", err)
	}
}
