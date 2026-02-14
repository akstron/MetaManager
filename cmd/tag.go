/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/heroku/self/MetaManager/internal/cmderror"
	"github.com/heroku/self/MetaManager/internal/data"
	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/filesys"
	"github.com/heroku/self/MetaManager/internal/storage"
	"github.com/heroku/self/MetaManager/internal/utils"

	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/spf13/cobra"
)

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Tagging related commands",
	Long:  "Tagging related commands",
}

func init() {
	rootCmd.AddCommand(tagCmd)

	// Register tag subcommands
	tagCmd.AddCommand(tagAddCmd)
	tagCmd.AddCommand(tagDeleteCmd)
	tagCmd.AddCommand(tagGetCmd)

	// Register node listTag command (tag-related)
	nodeCmd.AddCommand(nodeListTagCmd)

	// Register searchTag command (tag-related)
	searchCmd.AddCommand(searchTagCmd)
}

// tagAddInternal adds a tag to a file/directory
func tagAddInternal(ctxName string, args []string) error {
	rw, err := storage.GetRW(ctxName)
	if err != nil {
		return err
	}

	root, err := rw.Read()
	if err != nil {
		return err
	}

	tgMg := data.NewTagManager(data.NewDirTreeManager(ds.NewTreeManager(root)))

	resolver := filesys.NewBasicResolver(defaultStore)
	tagFilePath, err := resolver.Resolve(args[0])
	if err != nil {
		return err
	}

	tag := args[1]

	err = tgMg.AddTag(tagFilePath, tag)
	if err != nil {
		return err
	}

	err = tgMg.Save(rw)
	if err != nil {
		return err
	}

	return nil
}

func tagAdd(cmd *cobra.Command, args []string) {
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

	err = tagAddInternal(ctxName, args)
	if err != nil {
		goto finally
	}

finally:
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Tag added successfully")
	}
}

// tagAddCmd represents the tagAdd command
var tagAddCmd = &cobra.Command{
	Use:     "tagAdd",
	Short:   "Adds tag to a file/dir",
	Long:    `Adds tag to a file/dir`,
	Run:     tagAdd,
	Aliases: []string{"add"},
}

// tagDeleteInternal deletes a tag from a file/directory
func tagDeleteInternal(ctxName, path, tag string) error {
	rw, err := storage.GetRW(ctxName)
	if err != nil {
		return err
	}

	root, err := rw.Read()
	if err != nil {
		return err
	}

	tgMg := data.NewTagManager(data.NewDirTreeManager(ds.NewTreeManager(root)))

	resolver := filesys.NewBasicResolver(defaultStore)
	absPath, err := resolver.Resolve(path)
	if err != nil {
		return err
	}

	err = tgMg.DeleteTag(absPath, tag)
	if err != nil {
		return err
	}

	err = rw.Write(root)
	if err != nil {
		return err
	}

	fmt.Printf("tag %s deleted successfully\n", tag)

	return nil
}

func tagDelete(cmd *cobra.Command, args []string) {
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

	err = tagDeleteInternal(ctxName, args[0], args[1])
	if err != nil {
		goto finally
	}

finally:
	if err != nil {
		fmt.Println(err)
	}
}

// tagDeleteCmd represents the tagDelete command
var tagDeleteCmd = &cobra.Command{
	Use:     "tagDelete",
	Short:   "Deletes tag from a node (i.e., file/dir)",
	Long:    "Deletes tag from a node (i.e., file/dir)",
	Aliases: []string{"delete"},
	Run:     tagDelete,
}

// tagGetInternal gets all files/directories with a particular tag
func tagGetInternal(ctxName, tag string) ([]string, error) {
	rw, err := storage.GetRW(ctxName)
	if err != nil {
		return nil, err
	}

	root, err := rw.Read()
	if err != nil {
		return nil, err
	}

	tgMg := data.NewTagManager(data.NewDirTreeManager(ds.NewTreeManager(root)))

	paths, err := tgMg.GetTaggedNodes(tag)
	if err != nil {
		return nil, err
	}

	return paths, nil
}

func tagGet(_ *cobra.Command, args []string) {
	var err error
	var paths []string
	var pr list.Writer
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

	paths, err = tagGetInternal(ctxName, args[0])
	if err != nil {
		goto finally
	}

	pr = list.NewWriter()
	for _, path := range paths {
		pr.AppendItem(path)
	}
	pr.SetStyle(list.StyleDefault)
	fmt.Println(pr.Render())
finally:
	if err != nil {
		fmt.Println(err)
		// Print stack trace in case of error
		debug.PrintStack()
	}
}

// tagGetCmd represents the tagGet command
var tagGetCmd = &cobra.Command{
	Use:     "tagGet",
	Short:   "Gets files/dirs with a particular tag",
	Long:    `Gets files/dirs with a particular tag`,
	Run:     tagGet,
	Aliases: []string{"get"},
}

// nodeListInternal lists tags for a file/directory
func nodeListInternal(ctxName, path string) ([]string, error) {
	rw, err := storage.GetRW(ctxName)
	if err != nil {
		return nil, err
	}

	root, err := rw.Read()
	if err != nil {
		return nil, err
	}

	tgMg := data.NewTagManager(data.NewDirTreeManager(ds.NewTreeManager(root)))

	resolver := filesys.NewBasicResolver(defaultStore)
	absPath, err := resolver.Resolve(path)
	if err != nil {
		return nil, err
	}

	tags, err := tgMg.GetNodeTags(absPath)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func nodeList(cmd *cobra.Command, args []string) {
	var err error
	var tags []string
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

	tags, err = nodeListInternal(ctxName, args[0])
	if err != nil {
		goto finally
	}

	fmt.Println(tags)

finally:
	if err != nil {
		fmt.Println(err)
		// Print stack trace in case of error
		debug.PrintStack()
	}
}

// nodeListTagCmd represents the listTag command
var nodeListTagCmd = &cobra.Command{
	Use:     "listTag",
	Short:   "List tags of a file/dir",
	Long:    "List tags of a file/dir",
	Run:     nodeList,
	Aliases: []string{"lt", "tag"},
}

// searchTagCmd represents the searchTag command
var searchTagCmd = &cobra.Command{
	Use:   "searchTag",
	Short: "Find files/dir in the saved tree using regex",
	Long:  `Find files/dir in the saved tree using regex`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("To be implemented")
	},
}
