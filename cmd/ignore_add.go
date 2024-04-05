/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github/akstron/MetaManager/pkg/cmderror"
	"github/akstron/MetaManager/pkg/config"
	"github/akstron/MetaManager/pkg/utils"
	"path/filepath"

	"github.com/spf13/cobra"
)

func ignoreAdd(cmd *cobra.Command, args []string) {
	var err error
	var igMg *config.IgnoreManager
	pathToAdd := args[0]
	var absPathToAdd string

	isInitialized, err := utils.IsRootInitialized()
	if err != nil {
		goto finally
	}

	if !isInitialized {
		err = &cmderror.UninitializedRoot{}
		goto finally
	}

	igMg, err = config.NewIgnoreManager()
	if err != nil {
		goto finally
	}

	err = igMg.Load()
	if err != nil {
		goto finally
	}

	absPathToAdd, err = filepath.Abs(pathToAdd)
	if err != nil {
		goto finally
	}

	err = igMg.Add(absPathToAdd)
	if err != nil {
		goto finally
	}

	err = igMg.Save()
	if err != nil {
		goto finally
	}

finally:
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Path added to ignore list")
	}
}

// ignoreAddCmd represents the ignoreAdd command
var ignoreAddCommand = &cobra.Command{
	Use:   "ignoreAdd",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: ignoreAdd,
	/*
		This is used to provided aliases.
		i.e., instead of using ignoreAddCommand, we can also use add
	*/
	Aliases: []string{"add"},
}

func init() {
	ignoreCmd.AddCommand(ignoreAddCommand)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ignoreAddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ignoreAddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
