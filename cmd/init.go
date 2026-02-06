package cmd

import (
	"fmt"

	"github.com/heroku/self/MetaManager/internal/cmdmsg"
	"github.com/spf13/cobra"
)

// initCmd represents the init command.
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes the root of the directory which you want to manage",
	Long:  `Initializes the root of the directory which you want to manage.`,
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("%s the command expects exactly one argument", cmdmsg.ErrorOccurredMessage())
	}
	if err := InitRoot(args[0]); err != nil {
		return fmt.Errorf("%s: %w", cmdmsg.ErrorOccurredMessage(), err)
	}
	fmt.Println("Root path initialized successfully")
	return nil
}
