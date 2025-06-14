package cmd

import "github.com/spf13/cobra"

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Manage nodes",
}

func init() {
	nodeCmd.AddCommand(nodeAddCmd)
	rootCmd.AddCommand(nodeCmd)
}
