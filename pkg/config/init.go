package config

import (
	"encoding/json"
	"github/akstron/MetaManager/pkg/cmderror"
	"github/akstron/MetaManager/pkg/utils"
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

	const configDirName = `.mm`
	const configFileName = `config.json`
	const configFileIgnoreName = `ignore.json`
	const dataFileName = `data.json`

	dirPath, err := filepath.Abs(loc)
	if err != nil {
		return err
	}

	configDirPath := filepath.Join(dirPath, configDirName)

	/*
		If the current directory is already initialized which is indicated
		by the presence of .mm directory, we print an error accordingly.
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

	configFilePath := filepath.Join(configDirPath, configFileName)

	_, err = os.Create(configFilePath)
	if err != nil {
		return err
	}

	configFileIgnorePath := filepath.Join(configDirPath, configFileIgnoreName)

	_, err = os.Create(configFileIgnorePath)
	if err != nil {
		return err
	}

	dataFilePath := filepath.Join(configDirPath, dataFileName)
	_, err = os.Create(dataFilePath)
	if err != nil {
		return err
	}

	/*
		Write init info into config.json as json
	*/
	config := Config{RootPath: dirPath}
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFilePath, data, 0666)
	if err != nil {
		return err
	}

	return nil
}
