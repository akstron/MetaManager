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

func HandleSubtreeRemoval(pathExp string, drMg *data.DirTreeManager) error {
	if pathExp[len(pathExp)-1] == '*' {
		dirPath := pathExp[0 : len(pathExp)-1]
		dirPathAbs, err := filepath.Abs(dirPath)
		if err != nil {
			return err
		}

		err = drMg.SplitChildrenFromPath(dirPathAbs)
		if err != nil {
			return err
		}

		return nil
	}

	dirPathAbs, err := filepath.Abs(pathExp)
	if err != nil {
		return err
	}

	found, rootDirPathAbs, err := utils.FindRootDir()
	if err != nil {
		return err
	}

	if !found {
		return &cmderror.Unexpected{}
	}

	if dirPathAbs == rootDirPathAbs {
		return fmt.Errorf("utracking root folder is not allowed")
	}

	err = drMg.SplitNodeWithPath(dirPathAbs)
	if err != nil {
		return err
	}

	return nil
}

func untrackInternal(pathExp string) error {
	rw, err := storage.GetRW()
	if err != nil {
		return err
	}

	root, err := rw.Read()
	if err != nil {
		return err
	}

	drMg := data.NewDirTreeManager(ds.NewTreeManager(root))

	err = HandleSubtreeRemoval(pathExp, drMg)
	if err != nil {
		return err
	}

	err = rw.Write(drMg.Root)
	if err != nil {
		return err
	}

	return nil
}

func untrack(cmd *cobra.Command, args []string) {
	var err error

	if len(args) != 1 {
		err = &cmderror.InvalidNumberOfArguments{}
		goto finally
	}

	_, err = utils.CommonAlreadyInitializedChecks()
	if err != nil {
		goto finally
	}

	err = untrackInternal(args[0])
	if err != nil {
		goto finally
	}

finally:
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Location untracked successfully")
	}
}

// untrackCmd represents the untrack command
var untrackCmd = &cobra.Command{
	Use:   "untrack",
	Short: "Untracks an entire subtree rooted at node",
	Long:  "Untracks an entire subtree rooted at node",
	Run:   untrack,
}

func init() {
	rootCmd.AddCommand(untrackCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// untrackCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// untrackCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
