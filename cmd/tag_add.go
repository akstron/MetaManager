/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github/akstron/MetaManager/pkg/cmderror"
	"github/akstron/MetaManager/pkg/data"
	"github/akstron/MetaManager/pkg/utils"
	"path/filepath"
	"runtime/debug"

	"github.com/spf13/cobra"
)

func tagAddInternal(args []string) error {
	rw, err := GetRW()
	if err != nil {
		return err
	}

	tgMg := data.NewTagManager()

	err = tgMg.Load(rw)
	if err != nil {
		return err
	}

	tagFilePath, err := filepath.Abs(args[0])
	if err != nil {
		return err
	}
	tag := args[1]

	err = tgMg.AddTag(tagFilePath, tag)
	if err != nil {
		return err
	}

	err = tgMg.Save(rw)
	if err != nil {
		return err
	}

	return nil
}

func tagAdd(cmd *cobra.Command, args []string) {
	var err error

	if len(args) != 2 {
		err = &cmderror.InvalidNumberOfArguments{}
		goto finally
	}

	_, err = utils.CommonAlreadyInitializedChecks()
	if err != nil {
		goto finally
	}

	err = tagAddInternal(args)
	if err != nil {
		goto finally
	}

finally:
	if err != nil {
		fmt.Println(err)
		// Print stack trace in case of error
		debug.PrintStack()
	} else {
		fmt.Println("Tag added successfully")
	}
}

// tagAddCmd represents the tagAdd command
var tagAddCmd = &cobra.Command{
	Use:   "tagAdd",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run:     tagAdd,
	Aliases: []string{"add"},
}

func init() {
	tagCmd.AddCommand(tagAddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tagAddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tagAddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
