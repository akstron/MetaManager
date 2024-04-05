/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func ignoreList(cmd *cobra.Command, args []string) {
	fmt.Println("Implement later")
}

// ignoreListCmd represents the ignoreList command
var ignoreListCmd = &cobra.Command{
	Use:   "ignoreList",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run:     ignoreList,
	Aliases: []string{"list"},
}

func init() {
	ignoreCmd.AddCommand(ignoreListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ignoreListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ignoreListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
