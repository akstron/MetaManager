package config

import (
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
	const configFileName = `.mmconfig`

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

	return nil
}
