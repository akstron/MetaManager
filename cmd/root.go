// Package cmd provides the CLI commands for MetaManager.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd is the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "PathTracer",
	Short: "Manage your paths using this!",
	Long:  `Same as short description...`,
}

// Execute runs the root command and exits with the appropriate code on error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
