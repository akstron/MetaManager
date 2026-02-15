// Package cmd provides the CLI commands for MetaManager.
package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// RootCmd is the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "MetaManager",
	Short: "Manage your metadata using this!",
	Long:  `MetaManager is a tool for managing your files metadata.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if ok, _ := cmd.Root().PersistentFlags().GetBool("debug"); ok {
			logrus.SetLevel(logrus.DebugLevel)
		} else {
			logrus.SetLevel(logrus.InfoLevel)
		}
	},
}

func init() {
	RootCmd.PersistentFlags().BoolP("debug", "d", false, "enable debug logging")
}

// Execute runs the root command and exits with the appropriate code on error.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
