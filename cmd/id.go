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

// idCmd represents the id command
var idCmd = &cobra.Command{
	Use:   "id",
	Short: "id is an unique string which can be assigned to any node and used later for searching",
	Long:  "id is an unique string which can be assigned to any node and used later for searching",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("id called")
	},
}

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

// idSetCmd represents the set command
var idSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets id for a particular node",
	Long:  "Sets id for a particular node",
	Run:   idSet,
}

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

// idJumpCmd represents the jump command
var idJumpCmd = &cobra.Command{
	Use:   "jump",
	Short: "Jumps to the dir path or parent of a file",
	Long:  "Jumps to the dir path or parent of a file",
	Run:   idJump,
}

func getIdInternal(ctxName, path string) error {
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

	err = getIdInternal(ctxName, args[0])
	if err != nil {
		goto finally
	}

finally:
	if err != nil {
		fmt.Println(err)
	}
}

// idGetCmd represents the get command
var idGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets the id of the file/dir. Return <empty> if no id set",
	Long:  "Gets the id of the file/dir. Return <empty> if no id set",
	Run:   getId,
}

func init() {
	RootCmd.AddCommand(idCmd)

	idCmd.AddCommand(idSetCmd)
	idCmd.AddCommand(idJumpCmd)
	idCmd.AddCommand(idGetCmd)
}
