package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/heroku/self/MetaManager/internal/cmderror"
	"github.com/heroku/self/MetaManager/internal/config"
	dataPkg "github.com/heroku/self/MetaManager/internal/data"
	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/file"
	contextrepo "github.com/heroku/self/MetaManager/internal/repository/filesys"
	"github.com/heroku/self/MetaManager/internal/storage"
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
		Info:       &file.FileNode{GeneralNode: file.GeneralNode{AbsPath: absPath}},
		Children:   nil,
		Serializer: file.FileNodeJSONSerializer{},
	}
	dataFilePath := filepath.Join(appDir, utils.DataFileName)
	rw, err := storage.NewFileStorageRW(dataFilePath)
	if err != nil {
		return err
	}
	return rw.Write(emptyRoot)
}

/*
Initialize the root directory named '.mm', with a config file named '.mmconfig' in it,
for the current path provided.
*/
func InitRoot(loc string) error {
	if exist, err := utils.IsFilePresent(loc); err != nil {
		return err
	} else if !exist {
		return &cmderror.InvalidPath{}
	}

	dirPath, err := filepath.Abs(loc)
	if err != nil {
		return err
	}

	configDirPath := filepath.Join(dirPath, utils.MMDirName)

	/*
		If the current directory is already initialized which is indicated
		by the presence of .mm directory, we print an error accordingly.

		TODO: Add reinitialization with --force flag
	*/
	if exist, err := utils.IsFilePresent(configDirPath); err != nil {
		return err
	} else if exist {
		return &cmderror.AlreadyInitPath{}
	}

	err = os.Mkdir(configDirPath, 0755)
	if err != nil {
		return err
	}

	configFilePath := filepath.Join(configDirPath, utils.ConfigFileName)

	_, err = os.Create(configFilePath)
	if err != nil {
		return err
	}

	dataFilePath := filepath.Join(configDirPath, utils.DataFileName)
	_, err = os.Create(dataFilePath)
	if err != nil {
		return err
	}

	/*
		Write init info into config.json as json
	*/
	config := config.Config{RootPath: dirPath}
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFilePath, data, 0666)
	if err != nil {
		return err
	}

	rw, err := storage.NewFileStorageRW(dataFilePath)
	if err != nil {
		return err
	}

	mg := &dataPkg.DirTreeManager{}
	err = mg.MergeNodeWithPath(dirPath)
	if err != nil {
		return err
	}

	err = rw.Write(mg.Root)
	if err != nil {
		return err
	}

	return nil
}
