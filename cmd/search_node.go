/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/data"
	"github.com/heroku/self/MetaManager/internal/file"
	"github.com/heroku/self/MetaManager/internal/printer"
	"github.com/heroku/self/MetaManager/internal/utils"
	"github.com/heroku/self/MetaManager/internal/storage"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

func searchNode(cmd *cobra.Command, args []string) {
	var err error

	_, err = utils.CommonAlreadyInitializedChecks()
	if err != nil {
		goto finally
	}

	if len(args) != 1 {
		err = errors.New("this command accepts only 1 argument")
		goto finally
	}

	err = searchNodeInternal(args[0])
	if err != nil {
		goto finally
	}

finally:
	if err != nil {
		fmt.Println(err)
	}
}

func searchNodeInternal(regexPattern string) error {
	rw, err := storage.GetRW()
	if err != nil {
		return err
	}

	root, err := rw.Read()
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	drMg := data.NewDirTreeManager(ds.NewTreeManager(root))
	wdTrNode, err := drMg.FindTreeNodeByAbsPath(wd)
	if err != nil {
		return err
	}
	wdDrMg := data.NewDirTreeManager(ds.NewTreeManager(wdTrNode))

	foundTreeNodes, err := wdDrMg.FindTreeNodesByRegex(regexPattern)
	if err != nil {
		return err
	}

	drMgFound, err := data.BuildCopyTree(wd, foundTreeNodes)
	if err != nil {
		return err
	}

	pr := printer.NewTreePrinterManager(drMgFound.TreeManager)

	getPrintStringFunc := func(info any) (string, error) {
		node, ok := info.(file.NodeInformable)
		if !ok {
			return "", errors.New("info not convertible to NodeInformable")
		}

		str := node.GetAbsPath()
		match, err := regexp.MatchString(regexPattern, str)
		if err != nil {
			return "", err
		}

		str, err = utils.GetCurNodeFromAbsPath(str)
		if err != nil {
			return "", err
		}

		if match {
			str = str + " [found]"
		}

		return str, nil
	}
	prCxt := file.NewNodeExtraInfoPrinter(getPrintStringFunc)

	err = pr.TrPrintV2([]printer.PrintingContext{prCxt})
	if err != nil {
		return err
	}

	return nil
}

// searchNodeCmd represents the searchNode command
var searchNodeCmd = &cobra.Command{
	Use:   "searchNode",
	Short: "Find any file/directory in the saved tree using regex",
	Long: `Find any file/directory in the saved tree using regex. This prints the found nodes
in a tree fashion.`,
	Run:     searchNode,
	Aliases: []string{"node"},
}

func init() {
	searchCmd.AddCommand(searchNodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchNodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchNodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
