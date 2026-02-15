/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/heroku/self/MetaManager/internal/cmderror"
	"github.com/heroku/self/MetaManager/internal/data"
	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/repository/tree"
	"github.com/heroku/self/MetaManager/internal/utils"

	"github.com/spf13/cobra"
)

func idJumpInternal(ctxName, id string) error {
	rw, err := tree.GetRW(ctxName)
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

	err = idJumpInternal(ctxName, args[0])
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
