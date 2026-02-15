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
	"github.com/heroku/self/MetaManager/internal/file"
	"github.com/heroku/self/MetaManager/internal/filesys"
	"github.com/heroku/self/MetaManager/internal/printer"
	"github.com/heroku/self/MetaManager/internal/repository/tree"
	"github.com/heroku/self/MetaManager/internal/utils"
	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/sirupsen/logrus"
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
	tagCmd.AddCommand(searchTagCmd)
	tagCmd.AddCommand(tagListCmd)

	// Register flags for searchTag command
	searchTagCmd.Flags().BoolP("tree", "t", false, "Output results in tree format")
}

// tagAddInternal adds a tag to a file/directory
func tagAddInternal(ctxName string, args []string) error {
	rw, err := tree.GetRW(ctxName)
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
	rw, err := tree.GetRW(ctxName)
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
func tagSearchInternal(ctxName, tag string) ([]string, error) {
	rw, err := tree.GetRW(ctxName)
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

func tagSearch(cmd *cobra.Command, args []string) {
	var err error
	var paths []string
	var pr list.Writer
	var ctxName string
	var treeFlag bool

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

	treeFlag, err = cmd.Flags().GetBool("tree")
	if err != nil {
		goto finally
	}

	paths, err = tagSearchInternal(ctxName, args[0])
	if err != nil {
		goto finally
	}

	if treeFlag {
		logrus.Debugf("[tagSearch] treeFlag is true")
		err = tagSearchTreeInternal(ctxName, args[0], paths)
		if err != nil {
			goto finally
		}
		return
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

// tagSearchTreeInternal prints tagged nodes in tree format
func tagSearchTreeInternal(ctxName, tag string, paths []string) error {
	rw, err := tree.GetRW(ctxName)
	if err != nil {
		return err
	}

	root, err := rw.Read()
	if err != nil {
		return err
	}

	drMg := data.NewDirTreeManager(ds.NewTreeManager(root))

	// Get root path for building the tree
	rootInfo, ok := root.Info.(file.NodeInformable)
	if !ok {
		return fmt.Errorf("root info is not a NodeInformable")
	}
	rootPath := rootInfo.GetAbsPath()
	logrus.Debugf("[tagSearchTreeInternal] Root path: %s", rootPath)
	// Find tree nodes for each tagged path
	treeNodes := []*ds.TreeNode{}
	for _, path := range paths {
		node, err := drMg.FindTreeNodeByAbsPath(path)
		if err != nil {
			// Skip paths that don't exist in the tree
			continue
		}
		logrus.Debugf("[tagSearchTreeInternal] Found node: %s", node.Info.(file.NodeInformable).GetAbsPath())
		treeNodes = append(treeNodes, node)
	}

	if len(treeNodes) == 0 {
		fmt.Println("No tagged nodes found")
		return nil
	}

	// Build a tree from the found nodes
	drMgFound, err := data.BuildCopyTree(rootPath, treeNodes)
	if err != nil {
		return err
	}

	// Print the tree
	pr := printer.NewTreePrinterManager(drMgFound.TreeManager)
	err = pr.TrPrint([]string{"node"})
	if err != nil {
		return err
	}

	return nil
}

// searchTagCmd represents the searchTag command
var searchTagCmd = &cobra.Command{
	Use:     "searchTag",
	Short:   "Gets files/dirs with a particular tag",
	Long:    `Gets files/dirs with a particular tag. Use --tree flag to output results in tree format.`,
	Run:     tagSearch,
	Aliases: []string{"search"},
}

// tagGetInternal lists tags for a file/directory
func tagGetInternal(ctxName, path string) ([]string, error) {
	rw, err := tree.GetRW(ctxName)
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

func tagList(cmd *cobra.Command, args []string) {
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

	tags, err = tagGetInternal(ctxName, args[0])
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

// tagListCmd represents the list command (lists tags of a file/dir)
var tagListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List tags of a file/dir",
	Long:    "List tags of a file/dir",
	Run:     tagList,
	Aliases: []string{"ls"},
}
