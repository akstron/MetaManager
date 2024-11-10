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

func nodeListTracksInternal(tagFlag, idFlag bool) error {
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

	typesOfPrinting := []string{"node"}
	if idFlag {
		typesOfPrinting = append(typesOfPrinting, "id")
	}
	if tagFlag {
		typesOfPrinting = append(typesOfPrinting, "tags")
	}

	err = pr.TrPrint(typesOfPrinting)
	if err != nil {
		return err
	}

	return nil
}

// nodeListTrackCmd represents the nodeListTrack command
var nodeListTrackCmd = &cobra.Command{
	Use:     "nodeListTrack",
	Short:   "Lists all the tracked files/dirs from a particular root in a tree structure",
	Long:    "Lists all the tracked files/dirs from a particular root in a tree structure",
	Run:     nodeListTracks,
	Aliases: []string{"ltrack", "ltr", "tr", "tracks"},
}

func nodeListTracks(cmd *cobra.Command, args []string) {
	var err error
	var tagFlag, idFlag bool

	_, err = utils.CommonAlreadyInitializedChecks()
	if err != nil {
		goto finally
	}

	tagFlag, err = cmd.Flags().GetBool("tag")
	if err != nil {
		goto finally
	}

	idFlag, err = cmd.Flags().GetBool("id")
	if err != nil {
		goto finally
	}

	err = nodeListTracksInternal(tagFlag, idFlag)
	if err != nil {
		goto finally
	}

finally:
	if err != nil {
		fmt.Println(err)
	}
}

func init() {
	nodeCmd.AddCommand(nodeListTrackCmd)

	nodeListTrackCmd.PersistentFlags().BoolP("tag", "t", false, "flag to enrich the listing with tags for each node")
	nodeListTrackCmd.PersistentFlags().BoolP("id", "i", false, "flag to enrich the listing with id for each node")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodeListTrackCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodeListTrackCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
