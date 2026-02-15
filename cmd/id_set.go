/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/heroku/self/MetaManager/internal/cmderror"
	"github.com/heroku/self/MetaManager/internal/data"
	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/repository/tree"
	"github.com/heroku/self/MetaManager/internal/utils"

	"github.com/spf13/cobra"
)

func idSetInternal(ctxName, path, id string) error {
	idFilePath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	rw, err := tree.GetRW(ctxName)
	if err != nil {
		return err
	}

	root, err := rw.Read()
	if err != nil {
		return err
	}

	mg := data.NewDirTreeManager(ds.NewTreeManager(root))

	node, err := mg.FindFileNodeById(id)
	if err == nil {
		return fmt.Errorf("id: %s is already set for node %s", id, node.GetAbsPath())
	}

	pathNode, err := mg.FindNodeByAbsPath(idFilePath)
	if err != nil {
		return err
	}

	pathNode.SetId(id)

	err = rw.Write(mg.Root)
	if err != nil {
		return err
	}

	return nil
}

func idSet(cmd *cobra.Command, args []string) {
	var err error
	var ctxName string

	if len(args) != 2 {
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

	err = idSetInternal(ctxName, args[0], args[1])
	if err != nil {
		goto finally
	}

finally:
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("id set successfully")
	}
}

// setCmd represents the set command
var idSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets id for a particular node",
	Long:  "Sets id for a particular node",
	Run:   idSet,
}

func init() {
	idCmd.AddCommand(idSetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
