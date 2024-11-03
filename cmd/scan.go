/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github/akstron/MetaManager/filesys"
	"github/akstron/MetaManager/pkg/cmderror"
	"github/akstron/MetaManager/pkg/utils"
	"github/akstron/MetaManager/storage"
	"path/filepath"

	"github.com/spf13/cobra"
)

func scanInternal(rootDirPath string) error {
	rw, err := storage.GetRW()
	if err != nil {
		return err
	}

	rootNode, err := filesys.ScanDirectory(rootDirPath)
	if err != nil {
		return err
	}

	err = storage.WriteTree(rw, rootNode)
	if err != nil {
		return err
	}

	return nil
}

func scanPath(cmd *cobra.Command, args []string) {
	var err error
	var found bool
	var rootDirPath, mmDirPath string

	isInitialized, err := utils.IsRootInitialized()
	if err != nil {
		goto finally
	}

	if !isInitialized {
		err = &cmderror.UninitializedRoot{}
		goto finally
	}

	found, mmDirPath, err = utils.FindMMDirPath()
	if err != nil {
		goto finally
	}

	if !found {
		err = &cmderror.UninitializedRoot{}
		goto finally
	}

	rootDirPath = filepath.Join(mmDirPath, "../")

	err = scanInternal(rootDirPath)
	if err != nil {
		goto finally
	}

finally:
	if err != nil {
		fmt.Println(err)
	}
}

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scans the entire root directory for tracking --deprecate",
	Long:  "Scans the entire root directory for tracking --deprecate",
	Run:   scanPath,
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
