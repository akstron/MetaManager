/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github/akstron/MetaManager/ds"
	"github/akstron/MetaManager/pkg/data"
	"github/akstron/MetaManager/pkg/file"
	"github/akstron/MetaManager/pkg/utils"
	"github/akstron/MetaManager/storage"
	"os"
	"runtime/debug"
	"strings"

	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/spf13/cobra"
)

func nodeListTracksInternal() error {
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

	pr := list.NewWriter()

	err = ConstructWriter(requiredNode, "", pr)
	if err != nil {
		return err
	}

	pr.SetStyle(list.StyleConnectedLight)
	fmt.Println(pr.Render())

	return nil
}

func ConstructWriter(curNode *ds.TreeNode, cutPrefix string, wr list.Writer) error {
	info := curNode.Info.(file.NodeInformable)
	insPath, _ := strings.CutPrefix(info.GetAbsPath(), cutPrefix+"/")
	wr.AppendItem(insPath)

	wr.Indent()

	for _, child := range curNode.Children {
		err := ConstructWriter(child, info.GetAbsPath(), wr)
		if err != nil {
			return err
		}
	}

	wr.UnIndent()
	return nil
}

func nodeListTracks(cmd *cobra.Command, args []string) {
	var err error

	_, err = utils.CommonAlreadyInitializedChecks()
	if err != nil {
		goto finally
	}

	err = nodeListTracksInternal()
	if err != nil {
		goto finally
	}

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
