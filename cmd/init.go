/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github/akstron/MetaManager/pkg/cmdmsg"

	"github.com/spf13/cobra"
)

func initConfig(cmd *cobra.Command, locs []string) {
	if len(locs) != 1 {
		fmt.Println(cmdmsg.ErrorOccurredMessage(), " The command expects only 1 argument")
		return
	}

	err := InitRoot(locs[0])
	if err != nil {
		fmt.Println(cmdmsg.ErrorOccurredMessage(), err)
		return
	}

	fmt.Println("Root path initialized successfully")
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes the root of the directory which you want manage",
	Long:  `Same as short description`,
	Run:   initConfig,
}

func init() {
	rootCmd.AddCommand(initCmd)
}
