/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github/akstron/MetaManager/pkg/cmderror"
	"github/akstron/MetaManager/pkg/data"
	"github/akstron/MetaManager/pkg/utils"
	"path/filepath"

	"github.com/spf13/cobra"
)

/*
This will be changed based on certain flags -> Currently not implemented
*/
func GetRW() (data.TreeRW, error) {
	found, root, err := utils.FindMMDirPath()
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, &cmderror.UninitializedRoot{}
	}

	/*
		Check if the data.json is already written.
		Don't override, if already written
	*/
	dataFilePath, err := filepath.Abs(root + "/data.json")
	if err != nil {
		return nil, err
	}

	return data.NewFileStorageRW(dataFilePath)
}

func scanPath(cmd *cobra.Command, args []string) {
	var err error
	var root, rootDirPath string
	var mg data.TreeManager
	var rw data.TreeRW

	isInitialized, err := utils.IsRootInitialized()
	if err != nil {
		goto finally
	}

	if !isInitialized {
		err = &cmderror.UninitializedRoot{}
		goto finally
	}

	root, err = utils.GetAbsMMDirPath()
	if err != nil {
		goto finally
	}

	// Get parent directory
	// TODO: Update this to extract path from config.json
	rootDirPath, err = filepath.Abs(root + "/..")
	if err != nil {
		goto finally
	}

	/*
		With this factory method, this code path is free from
		RW related changes.
	*/
	rw, err = GetRW()
	if err != nil {
		goto finally
	}

	// Do other heavy lifting only when data file is empty
	mg = data.TreeManager{
		// Parent of root which is .mm directory path
		DirPath: rootDirPath,
	}

	err = mg.ScanDirectory()
	if err != nil {
		goto finally
	}

	err = data.WriteTree(rw, mg.Root)
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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: scanPath,
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
