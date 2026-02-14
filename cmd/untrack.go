/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/heroku/self/MetaManager/internal/cmderror"
	"github.com/heroku/self/MetaManager/internal/data"
	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/filesys"
	"github.com/heroku/self/MetaManager/internal/storage"
	"github.com/heroku/self/MetaManager/internal/utils"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// pathExp would be an absolute path like "/Folder/SubFolder" or "gdrive:/Folder/SubFolder".
func HandleSubtreeRemoval(ctxName, pathExp string, drMg *data.DirTreeManager) error {
	if pathExp[len(pathExp)-1] == '*' {
		dirPath := pathExp[0 : len(pathExp)-1]

		err := drMg.SplitChildrenFromPath(dirPath)
		if err != nil {
			return err
		}

		return nil
	}

	found, rootDirPathAbs, err := utils.FindRootDir(ctxName)
	if err != nil {
		return err
	}

	if !found {
		return &cmderror.Unexpected{}
	}

	if pathExp == rootDirPathAbs {
		return fmt.Errorf("utracking root folder is not allowed")
	}

	err = drMg.SplitNodeWithPath(pathExp)
	if err != nil {
		return err
	}

	return nil
}

func untrackInternal(ctxName, pathExp string) error {
	rw, err := storage.GetRW(ctxName)
	if err != nil {
		return err
	}

	root, err := rw.Read()
	if err != nil {
		return err
	}

	drMg := data.NewDirTreeManager(ds.NewTreeManager(root))

	resolver := filesys.NewBasicResolver(defaultStore)
	resolvedPath, err := resolver.Resolve(pathExp)
	if err != nil {
		return err
	}

	logrus.Debugf("[untrack] resolvedPath: %q", resolvedPath)

	err = HandleSubtreeRemoval(ctxName, resolvedPath, drMg)
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

	err = untrackInternal(ctxName, args[0])
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
