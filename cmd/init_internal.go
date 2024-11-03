package cmd

import (
	"encoding/json"
	"github/akstron/MetaManager/pkg/cmderror"
	"github/akstron/MetaManager/pkg/config"
	dataPkg "github/akstron/MetaManager/pkg/data"
	"github/akstron/MetaManager/pkg/utils"
	"github/akstron/MetaManager/storage"
	"os"
	"path/filepath"
)

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

	configDirPath := filepath.Join(dirPath, utils.MM_DIR_NAME)

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

	configFilePath := filepath.Join(configDirPath, utils.CONFIG_FILE_NAME)

	_, err = os.Create(configFilePath)
	if err != nil {
		return err
	}

	configFileIgnorePath := filepath.Join(configDirPath, utils.IGNORE_FILE_NAME)

	_, err = os.Create(configFileIgnorePath)
	if err != nil {
		return err
	}

	dataFilePath := filepath.Join(configDirPath, utils.DATA_FILE_NAME)
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

	rw, err := storage.GetRW()
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
