package config

import (
	"encoding/json"
	"github.com/heroku/self/MetaManager/internal/utils"
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

// NewIgnoreManager creates an IgnoreManager for the given .mm directory path (e.g. .mm/<contextName>/).
// If mmDirPath is empty, returns a manager with empty ignore list (no file path).
func NewIgnoreManager(mmDirPath string) (*IgnoreManager, error) {
	ig := &IgnoreManager{}
	if mmDirPath == "" {
		return ig, nil
	}
	ignoreFilePath, err := filepath.Abs(filepath.Join(mmDirPath, utils.IgnoreFileName))
	if err != nil {
		return nil, err
	}
	ig.ignoreFilePath = ignoreFilePath
	return ig, nil
}

func (mg *IgnoreManager) Save() error {
	if mg.ignoreFilePath == "" {
		return nil
	}
	content, err := json.Marshal(mg.Data)
	if err != nil {
		return err
	}
	return os.WriteFile(mg.ignoreFilePath, content, 0666)
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
	if mg.ignoreFilePath == "" {
		return nil
	}
	content, err := os.ReadFile(mg.ignoreFilePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, &mg.Data)
}
