package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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

func IsRootInitialized() (bool, error) {
	dirPath := "./.mm"
	return IsFilePresent(dirPath)
}

/*
This gives the root path of .mm directory
TODO: Change and test this. Use FindMMDirPath to get path
*/
func GetAbsMMDirPath() (string, error) {
	dirPath := "./.mm"
	return filepath.Abs(dirPath)
}

func IsFileEmpty(filePath string) (bool, error) {
	stat, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}

	return stat.Size() == 0, nil
}

func FindRootDir() (bool, string, error) {
	found, mmDirPath, err := FindMMDirPath()
	if err != nil {
		return false, "", err
	}

	if !found {
		return false, "", nil
	}

	return true, filepath.Join(mmDirPath, ".."), nil
}

func FindMMDirPath() (bool, string, error) {
	testEnvDir := os.Getenv("MM_TEST_ENV_DIR")
	wd, err := os.Getwd()
	if testEnvDir != "" {
		wd = testEnvDir
	}
	if err != nil {
		return false, "", err
	}
	return findMMDirPathInternal(wd)
}

/*
Recursively searches the parents for .mm directory
*/
func findMMDirPathInternal(path string) (bool, string, error) {
	for {
		isPresent, err := IsMMDirPresent(path)
		if err != nil {
			return false, "", err
		}

		if isPresent {
			mmDirPath, err := filepath.Abs(path + "/" + MMDirName)
			if err != nil {
				return false, "", err
			}
			return true, mmDirPath, nil
		}

		newPath := filepath.Dir(path)
		if path == newPath {
			break
		}
		path = newPath
	}

	return false, "", nil
}

func IsMMDirPresent(path string) (bool, error) {
	mmDirPath, err := filepath.Abs(path + "/" + MMDirName)
	if err != nil {
		return false, err
	}
	return IsFilePresent(mmDirPath)
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
