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

	"github.com/spf13/cobra"
)

func idJumpInternal(id string) error {
	rw, err := storage.GetRW()
	if err != nil {
		return err
	}

	root, err := rw.Read()
	if err != nil {
		return err
	}

	mg := data.NewDirTreeManager(ds.NewTreeManager(root))
	fileNode, err := mg.FindFileNodeById(id)
	if err != nil {
		return err
	}

	fmt.Println(fileNode.GetAbsPath())

	return nil
}

func idJump(cmd *cobra.Command, args []string) {
	var err error

	if len(args) != 1 {
		err = &cmderror.InvalidNumberOfArguments{}
		goto finally
	}

	_, err = utils.CommonAlreadyInitializedChecks()
	if err != nil {
		goto finally
	}

	err = idJumpInternal(args[0])
	if err != nil {
		goto finally
	}

finally:
	if err != nil {
		fmt.Println(err)
	}
}

// jumpCmd represents the jump command
var idJumpCmd = &cobra.Command{
	Use:   "jump",
	Short: "Jumps to the dir path or parent of a file",
	Long:  "Jumps to the dir path or parent of a file",
	Run:   idJump,
}

func init() {
	idCmd.AddCommand(idJumpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jumpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jumpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
