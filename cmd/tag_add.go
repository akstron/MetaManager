/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/cmderror"
	"github.com/heroku/self/MetaManager/internal/data"
	"github.com/heroku/self/MetaManager/internal/utils"
	"github.com/heroku/self/MetaManager/internal/storage"
	"path/filepath"

	"github.com/spf13/cobra"
)

func tagAddInternal(args []string) error {
	rw, err := storage.GetRW()
	if err != nil {
		return err
	}

	root, err := rw.Read()
	if err != nil {
		return err
	}

	tgMg := data.NewTagManager(data.NewDirTreeManager(ds.NewTreeManager(root)))

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
	} else {
		fmt.Println("Tag added successfully")
	}
}

// tagAddCmd represents the tagAdd command
var tagAddCmd = &cobra.Command{
	Use:     "tagAdd",
	Short:   "Adds tag to a file/dir",
	Long:    `Adds tag to a file/dir`,
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
