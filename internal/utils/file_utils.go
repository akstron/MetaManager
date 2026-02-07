package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetBaseDir returns the directory for app data (context file, .mm, etc.).
// Uses MM_TEST_CONTEXT_DIR in tests, otherwise the executable's directory.
func GetBaseDir() (string, error) {
	if d := os.Getenv("MM_TEST_CONTEXT_DIR"); d != "" {
		return d, nil
	}
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("executable path: %w", err)
	}
	return filepath.Dir(execPath), nil
}

// GetAppDataDir returns the path to the shared .mm parent directory (next to the executable).
// Per-context data lives under .mm/<contextName>/.
func GetAppDataDir() (string, error) {
	base, err := GetBaseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, MMDirName), nil
}

// GetAppDataDirForContext returns the path to the .mm directory for the given context (e.g. .mm/mydrive/).
// The directory may not exist until the context is created.
func GetAppDataDirForContext(contextName string) (string, error) {
	if contextName == "" {
		return "", fmt.Errorf("context name cannot be empty")
	}
	parent, err := GetAppDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(parent, contextName), nil
}

/*
All the mentioned utils can also be used for directories,
as dirs are more or less dirs
*/
func IsFilePresent(dirPath string) (bool, error) {
	_, err := os.Stat(dirPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// IsRootInitialized returns true if the .mm directory for the given context exists.
func IsRootInitialized(contextName string) (bool, error) {
	if contextName == "" {
		return false, nil
	}
	appDir, err := GetAppDataDirForContext(contextName)
	if err != nil {
		return false, err
	}
	return IsFilePresent(appDir)
}

// GetAbsMMDirPath returns the absolute path to the .mm directory for the given context.
func GetAbsMMDirPath(contextName string) (string, error) {
	return GetAppDataDirForContext(contextName)
}

func IsFileEmpty(filePath string) (bool, error) {
	stat, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}

	return stat.Size() == 0, nil
}

// FindRootDir returns the parent of the .mm dir for the context (the "root" path). Used for path checks.
func FindRootDir(contextName string) (bool, string, error) {
	found, mmDirPath, err := FindMMDirPath(contextName)
	if err != nil {
		return false, "", err
	}
	if !found {
		return false, "", nil
	}
	return true, filepath.Join(mmDirPath, ".."), nil
}

// FindMMDirPath returns the .mm directory path for the given context (.mm/<contextName>/).
// If contextName is empty or the directory does not exist, returns (false, "", nil).
func FindMMDirPath(contextName string) (bool, string, error) {
	if contextName == "" {
		return false, "", nil
	}
	appDir, err := GetAppDataDirForContext(contextName)
	if err != nil {
		return false, "", err
	}
	exists, err := IsFilePresent(appDir)
	if err != nil || !exists {
		return false, "", err
	}
	return true, appDir, nil
}

func SaveToFile(location string, data []byte) error {
	return os.WriteFile(location, data, 0666)
}

func GetCurNodeFromAbsPath(absPath string) (string, error) {
	temp := strings.Split(absPath, "/")
	if len(temp) == 0 {
		return "", fmt.Errorf("can't print path: %s", absPath)
	}
	return temp[len(temp)-1], nil
}
