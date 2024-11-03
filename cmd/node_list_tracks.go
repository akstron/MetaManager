/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github/akstron/MetaManager/pkg/utils"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// func nodeListTracksInternal() error {
// 	rw, err := storage.GetRW()
// 	if err != nil {
// 		return err
// 	}

// 	root, err := rw.Read()
// 	if err != nil {
// 		return err
// 	}

// 	dirPath, err := os.Getwd()
// 	if err != nil {
// 		return err
// 	}

// 	drMg := data.NewDirTreeManager(ds.NewTreeManager(root))

// 	requiredNode, err := drMg.FindNodeByAbsPath(dirPath)
// 	if err != nil {
// 		return err
// 	}

// 	paths := []string{}

// 	pr := list.NewWriter()

// 	iter := ds.NewTreeIterator(ds.NewTreeManager(requiredNode))
// 	for iter.HasNext() {
// 		got, err := iter.Next()
// 		if err != nil {
// 			return err
// 		}

// 		info, ok := got.(file.NodeInformable)
// 		if !ok {
// 			return &cmderror.Unexpected{}
// 		}
// 		paths = append(paths, info.GetAbsPath())
// 	}

// 	for _, path := range paths {
// 		pr.AppendItem(path)
// 	}
// 	pr.SetStyle(list.StyleDefault)
// 	fmt.Println(pr.Render())

// 	return nil
// }

func nodeListTracks(cmd *cobra.Command, args []string) {
	var err error

	_, err = utils.CommonAlreadyInitializedChecks()
	if err != nil {
		goto finally
	}

	// err = nodeListTracksInternal()
	// if err != nil {
	// 	goto finally
	// }

finally:
	if err != nil {
		fmt.Println(err)
		// Print stack trace in case of error
		debug.PrintStack()
	}
}

// nodeListTrackCmd represents the nodeListTrack command
var nodeListTrackCmd = &cobra.Command{
	Use:     "nodeListTrack",
	Short:   "Lists all the tracked files/dirs from a particular root in a tree structure",
	Long:    "Lists all the tracked files/dirs from a particular root in a tree structure",
	Run:     nodeListTracks,
	Aliases: []string{"ltrack"},
}

func init() {
	nodeCmd.AddCommand(nodeListTrackCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodeListTrackCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodeListTrackCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
