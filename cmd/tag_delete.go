/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github/akstron/MetaManager/ds"
	"github/akstron/MetaManager/pkg/cmderror"
	"github/akstron/MetaManager/pkg/data"
	"github/akstron/MetaManager/pkg/utils"
	"github/akstron/MetaManager/storage"
	"path/filepath"

	"github.com/spf13/cobra"
)

func tagDeleteInternal(path, tag string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	rw, err := storage.GetRW()
	if err != nil {
		return err
	}

	root, err := rw.Read()
	if err != nil {
		return err
	}

	tgMg := data.NewTagManager(data.NewDirTreeManager(ds.NewTreeManager(root)))

	err = tgMg.DeleteTag(absPath, tag)
	if err != nil {
		return err
	}

	err = rw.Write(root)
	if err != nil {
		return err
	}

	fmt.Printf("tag %s deleted successfully\n", tag)

	return nil
}

func tagDelete(cmd *cobra.Command, args []string) {
	var err error

	if len(args) != 2 {
		err = &cmderror.InvalidNumberOfArguments{}
		goto finally
	}

	_, err = utils.CommonAlreadyInitializedChecks()
	if err != nil {
		goto finally
	}

	err = tagDeleteInternal(args[0], args[1])
	if err != nil {
		goto finally
	}

finally:
	if err != nil {
		fmt.Println(err)
	}
}

// tagDeleteCmd represents the tagDelete command
var tagDeleteCmd = &cobra.Command{
	Use:     "tagDelete",
	Short:   "Deletes tag from a node (i.e., file/dir)",
	Long:    "Deletes tag from a node (i.e., file/dir)",
	Aliases: []string{"delete"},
	Run:     tagDelete,
}

func init() {
	tagCmd.AddCommand(tagDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tagDeleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tagDeleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
