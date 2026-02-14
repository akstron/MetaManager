/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/cmderror"
	"github.com/heroku/self/MetaManager/internal/data"
	"github.com/heroku/self/MetaManager/internal/utils"
	"github.com/heroku/self/MetaManager/internal/storage"
	"github.com/heroku/self/MetaManager/internal/file"
	contextrepo "github.com/heroku/self/MetaManager/internal/repository/filesys"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

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

	// Get context type to determine path handling
	ctxType, err := GetContextType(ctxName)
	if err != nil {
		return err
	}

	var absPath string
	if ctxType == contextrepo.TypeGDrive {
		// For gdrive context, resolve relative paths and prefix with gdrive:/
		cwd, err := defaultStore.GetGDriveCwd()
		if err != nil {
			return err
		}
		resolved := contextrepo.ResolveGDrivePath(cwd, path)
		// Remove leading slash and prefix with gdrive:/
		if resolved == "/" {
			absPath = file.GDrivePathPrefix
		} else {
			absPath = file.GDrivePathPrefix + strings.TrimPrefix(resolved, "/")
		}
	} else {
		// For local context, use absolute path
		absPath, err = filepath.Abs(path)
		if err != nil {
			return err
		}
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

func init() {
	tagCmd.AddCommand(tagDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tagDeleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tagDeleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
