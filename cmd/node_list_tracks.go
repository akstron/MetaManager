/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github/akstron/MetaManager/ds"
	"github/akstron/MetaManager/pkg/data"
	"github/akstron/MetaManager/pkg/printer"
	"github/akstron/MetaManager/pkg/utils"
	"github/akstron/MetaManager/storage"
	"os"

	"github.com/spf13/cobra"
)

func nodeListTracksInternal(tagFlag bool) error {
	rw, err := storage.GetRW()
	if err != nil {
		return err
	}

	root, err := rw.Read()
	if err != nil {
		return err
	}

	dirPath, err := os.Getwd()
	if err != nil {
		return err
	}

	drMg := data.NewDirTreeManager(ds.NewTreeManager(root))

	requiredNode, err := drMg.FindTreeNodeByAbsPath(dirPath)
	if err != nil {
		return err
	}

	pr := printer.NewTreePrinterManager(ds.NewTreeManager(requiredNode))

	typeOfPrinting := "node"
	if tagFlag {
		typeOfPrinting = "node-tags"
	}

	err = pr.TrPrint(typeOfPrinting)
	if err != nil {
		return err
	}

	return nil
}

func nodeListTracks(cmd *cobra.Command, args []string) {
	var err error
	// var tagFlag bool

	_, err = utils.CommonAlreadyInitializedChecks()
	if err != nil {
		goto finally
	}

	// tagFlag, err = cmd.Flags().GetBool("tag")
	// if err != nil {
	// 	goto finally
	// }

	err = nodeListTracksInternal(true)
	if err != nil {
		goto finally
	}

finally:
	if err != nil {
		fmt.Println(err)
	}
}

// nodeListTrackCmd represents the nodeListTrack command
var nodeListTrackCmd = &cobra.Command{
	Use:     "nodeListTrack",
	Short:   "Lists all the tracked files/dirs from a particular root in a tree structure",
	Long:    "Lists all the tracked files/dirs from a particular root in a tree structure",
	Run:     nodeListTracks,
	Aliases: []string{"ltrack", "ltr"},
}

func init() {
	nodeCmd.AddCommand(nodeListTrackCmd)

	nodeListTagCmd.PersistentFlags().Bool("tag", false, "Use --tag flag to enrich the listing with tags for each node")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodeListTrackCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodeListTrackCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
