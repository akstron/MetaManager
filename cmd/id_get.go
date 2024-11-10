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

func getIdInternal(path string) error {
	idFilePath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	rw, err := storage.GetRW()
	if err != nil {
		return err
	}

	root, err := rw.Read()
	if err != nil {
		return err
	}

	mg := data.NewDirTreeManager(ds.NewTreeManager(root))

	pathNode, err := mg.FindNodeByAbsPath(idFilePath)
	if err != nil {
		return err
	}

	id := pathNode.GetId()
	if id == "" {
		id = "<empty>"
	}

	fmt.Println(id)

	return nil
}

func getId(cmd *cobra.Command, args []string) {
	var err error

	if len(args) != 1 {
		err = &cmderror.InvalidNumberOfArguments{}
		goto finally
	}

	_, err = utils.CommonAlreadyInitializedChecks()
	if err != nil {
		goto finally
	}

	err = getIdInternal(args[0])
	if err != nil {
		goto finally
	}

finally:
	if err != nil {
		fmt.Println(err)
	}
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets the id of the file/dir. Return <empty> if no id set",
	Long:  "Gets the id of the file/dir. Return <empty> if no id set",
	Run:   getId,
}

func init() {
	idCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
