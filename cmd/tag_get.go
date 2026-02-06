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
	"runtime/debug"

	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/spf13/cobra"
)

func tagGetInternal(tag string) ([]string, error) {
	rw, err := storage.GetRW()
	if err != nil {
		return nil, err
	}

	root, err := rw.Read()
	if err != nil {
		return nil, err
	}

	tgMg := data.NewTagManager(data.NewDirTreeManager(ds.NewTreeManager(root)))

	paths, err := tgMg.GetTaggedNodes(tag)
	if err != nil {
		return nil, err
	}

	return paths, nil
}

func tagGet(_ *cobra.Command, args []string) {
	var err error
	var paths []string
	var pr list.Writer

	if len(args) != 1 {
		err = &cmderror.InvalidNumberOfArguments{}
		goto finally
	}

	_, err = utils.CommonAlreadyInitializedChecks()
	if err != nil {
		goto finally
	}

	paths, err = tagGetInternal(args[0])
	if err != nil {
		goto finally
	}

	pr = list.NewWriter()
	for _, path := range paths {
		pr.AppendItem(path)
	}
	pr.SetStyle(list.StyleDefault)
	fmt.Println(pr.Render())
finally:
	if err != nil {
		fmt.Println(err)
		// Print stack trace in case of error
		debug.PrintStack()
	}
}

// tagAddCmd represents the tagAdd command
var tagGetCmd = &cobra.Command{
	Use:     "tagGet",
	Short:   "Gets files/dirs with a particular tag",
	Long:    `Gets files/dirs with a particular tag`,
	Run:     tagGet,
	Aliases: []string{"get"},
}

func init() {
	tagCmd.AddCommand(tagGetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tagAddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tagAddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
