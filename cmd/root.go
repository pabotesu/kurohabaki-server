package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var (
    configPath string // Holds --config flag value
)

// rootCmd is the base command
var rootCmd = &cobra.Command{
    Use:   "kurohabaki-server",
    Short: "kurohabaki-server is a WireGuard-based P2P coordination server",
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        fmt.Println("Using config file:", configPath)
    },
}

// Execute runs the CLI entrypoint
func Execute() {
    rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to the configuration file (optional)")
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

// GetConfigPath exposes the parsed configPath to other packages
func GetConfigPath() string {
    return configPath
}
