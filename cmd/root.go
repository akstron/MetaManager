// Package cmd provides the CLI commands for MetaManager.
package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd is the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "PathTracer",
	Short: "Manage your paths using this!",
	Long:  `Same as short description...`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if ok, _ := cmd.Root().PersistentFlags().GetBool("debug"); ok {
			logrus.SetLevel(logrus.DebugLevel)
		} else {
			logrus.SetLevel(logrus.InfoLevel)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "enable debug logging")
}

// Execute runs the root command and exits with the appropriate code on error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
