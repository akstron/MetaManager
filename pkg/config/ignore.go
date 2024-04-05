package config

import (
	"encoding/json"
	"github/akstron/MetaManager/pkg/utils"
	"os"
	"path/filepath"
)

type IgnoreData struct {
	Paths []string
}

type IgnoreManager struct {
	Data           IgnoreData
	ignoreFilePath string
}

// Shift this to utils
func generateIgnoreFilePath() (string, error) {
	rootPath, err := utils.GetAbsMMDirPath()
	if err != nil {
		return "", err
	}

	ignoreFilePath, err := filepath.Abs(rootPath + "/" + utils.IGNORE_FILE_NAME)
	if err != nil {
		return "", err
	}

	return ignoreFilePath, nil
}

func NewIgnoreManager() (*IgnoreManager, error) {
	var err error
	ig := &IgnoreManager{}
	ig.ignoreFilePath, err = generateIgnoreFilePath()
	if err != nil {
		return nil, err
	}
	return ig, nil
}

func (mg *IgnoreManager) Save() error {
	content, err := json.Marshal(mg.Data)
	if err != nil {
		return err
	}

	err = os.WriteFile(mg.ignoreFilePath, content, 0666)
	if err != nil {
		return err
	}

	return nil
}

func (mg *IgnoreManager) Add(path string) error {
	mg.Data.Paths = append(mg.Data.Paths, path)
	return nil
}

/*
Implement caching here - Don't load again, if already loaded

Can make ForceLoad(), if we deliberately want to load
*/
func (mg *IgnoreManager) Load() error {
	content, err := os.ReadFile(mg.ignoreFilePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &mg.Data)
	if err != nil {
		return err
	}

	return nil
}
