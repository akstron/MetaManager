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

/*
This will be changed based on certain flags -> Currently not implemented
*/
func GetRW() (storage.TreeRW, error) {
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
	dataFilePath := filepath.Join(root, utils.DATA_FILE_NAME)

	return storage.NewFileStorageRW(dataFilePath)
}

func scanInternal(rootDirPath string) error {
	rw, err := GetRW()
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
