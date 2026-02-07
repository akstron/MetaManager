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
	"runtime/debug"

	"github.com/spf13/cobra"
)

func nodeListInternal(ctxName, path string) ([]string, error) {
	rw, err := storage.GetRW(ctxName)
	if err != nil {
		return nil, err
	}

	root, err := rw.Read()
	if err != nil {
		return nil, err
	}

	tgMg := data.NewTagManager(data.NewDirTreeManager(ds.NewTreeManager(root)))

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	tags, err := tgMg.GetNodeTags(absPath)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func nodeList(cmd *cobra.Command, args []string) {
	var err error
	var tags []string
	var ctxName string

	if len(args) != 1 {
		err = &cmderror.InvalidNumberOfArguments{}
		goto finally
	}

	ctxName, err = getContextRequired()
	if err != nil {
		goto finally
	}
	_, err = utils.CommonAlreadyInitializedChecks(ctxName)
	if err != nil {
		goto finally
	}

	tags, err = nodeListInternal(ctxName, args[0])
	if err != nil {
		goto finally
	}

	fmt.Println(tags)

finally:
	if err != nil {
		fmt.Println(err)
		// Print stack trace in case of error
		debug.PrintStack()
	}
}

// listTagCmd represents the listTag command
var nodeListTagCmd = &cobra.Command{
	Use:     "listTag",
	Short:   "List tags of a file/dir",
	Long:    "List tags of a file/dir",
	Run:     nodeList,
	Aliases: []string{"lt", "tag"},
}

func init() {
	nodeCmd.AddCommand(nodeListTagCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listTagCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listTagCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
