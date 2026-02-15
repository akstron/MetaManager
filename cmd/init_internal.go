package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/heroku/self/MetaManager/internal/config"
	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/file"
	contextrepo "github.com/heroku/self/MetaManager/internal/repository/context"
	"github.com/heroku/self/MetaManager/internal/repository/tree"
	"github.com/heroku/self/MetaManager/internal/utils"
)

// EnsureAppDataDir creates the .mm/<contextName> directory for the given context (next to the executable
// or under MM_TEST_CONTEXT_DIR) with config.json and data.json if it does not exist. Idempotent.
func EnsureAppDataDir(contextName string) error {
	if contextName == "" {
		return fmt.Errorf("context name cannot be empty")
	}
	parentDir, err := utils.GetAppDataDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return err
	}
	appDir, err := utils.GetAppDataDirForContext(contextName)
	if err != nil {
		return err
	}
	exists, err := utils.IsFilePresent(appDir)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	if err := os.Mkdir(appDir, 0755); err != nil {
		return err
	}
	baseDir, err := utils.GetBaseDir()
	if err != nil {
		return err
	}

	configFilePath := filepath.Join(appDir, utils.ConfigFileName)
	cfg := config.Config{RootPath: baseDir}
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	if err := os.WriteFile(configFilePath, data, 0666); err != nil {
		return err
	}

	absPath := "/"
	contextType, err := GetContextType(contextName)
	if err == nil && contextType == contextrepo.TypeGDrive {
		absPath = file.GDrivePathPrefix
	} else if os.Getenv("MM_TEST_CONTEXT_DIR") != "" {
		// In tests, root must match the track root so merge creates the right number of nodes.
		absPath = baseDir
	}
	emptyRoot := &ds.TreeNode{
		// root is a special path that is used to represent the root of the tree
		Info:     &file.FileNode{GeneralNode: file.GeneralNode{AbsPath: absPath}},
		Children: nil,
	}
	dataFilePath := filepath.Join(appDir, utils.DataFileName)
	rw, err := tree.NewFileStorageRW(dataFilePath)
	if err != nil {
		return err
	}
	return rw.Write(emptyRoot)
}
