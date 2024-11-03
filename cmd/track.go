/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github/akstron/MetaManager/ds"
	"github/akstron/MetaManager/filesys"
	"github/akstron/MetaManager/pkg/cmderror"
	"github/akstron/MetaManager/pkg/data"
	"github/akstron/MetaManager/pkg/utils"
	"github/akstron/MetaManager/storage"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"
)

func trackInternal(pathExp string) error {
	found, rootDirPath, err := utils.FindRootDir()
	if err != nil {
		return err
	}

	if !found {
		return &cmderror.Unexpected{}
	}

	isPathExpInRootDir := strings.HasPrefix(pathExp, rootDirPath)
	if !isPathExpInRootDir {
		return &cmderror.InvalidOperation{}
	}

	rw, err := storage.GetRW()
	if err != nil {
		return err
	}

	root, err := rw.Read()
	if err != nil {
		return err
	}

	subTree, err := filesys.Track(pathExp)
	if err != nil {
		return err
	}

	drMg := data.NewDirTreeManager(ds.NewTreeManager(root))

	/*
		Why no use MergeNodeWithPath ?
		This is because pathExp can contain * at the end
	*/
	drMg.MergeNode(subTree)

	err = rw.Write(drMg.Root)
	if err != nil {
		return err
	}

	return nil
}

func track(cmd *cobra.Command, args []string) {
	var err error

	if len(args) != 1 {
		err = &cmderror.InvalidNumberOfArguments{}
		goto finally
	}

	_, err = utils.CommonAlreadyInitializedChecks()
	if err != nil {
		goto finally
	}

	err = trackInternal(args[0])
	if err != nil {
		goto finally
	}

finally:
	if err != nil {
		fmt.Println(err)
		// Print stack trace in case of error
		debug.PrintStack()
	} else {
		fmt.Println("Location tracked successfully")
	}
}

// trackCmd represents the track command
var trackCmd = &cobra.Command{
	Use:   "track",
	Short: "Starts tracking a file/folder/all-files-and-folder-at-a-loc",
	Long:  "Starts tracking a file/folder/all-files-and-folder-at-a-loc",
	Run:   track,
}

func init() {
	rootCmd.AddCommand(trackCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// trackCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// trackCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
