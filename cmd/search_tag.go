/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// searchTagCmd represents the searchTag command
var searchTagCmd = &cobra.Command{
	Use:   "searchTag",
	Short: "Find files/dir in the saved tree using regex",
	Long:  `Find files/dir in the saved tree using regex`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("To be implemented")
	},
}

func init() {
	searchCmd.AddCommand(searchTagCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchTagCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchTagCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
