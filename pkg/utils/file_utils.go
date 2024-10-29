package utils

import (
	"os"
	"path/filepath"
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

/*
Not same as C++
we can't have same name as FindMMDirPath
TODO: Check this later
*/
func FindMMDirPath() (bool, string, error) {
	wd, err := os.Getwd()
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
			mmDirPath, err := filepath.Abs(path + "/" + MM_DIR_NAME)
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
	mmDirPath, err := filepath.Abs(path + "/" + MM_DIR_NAME)
	if err != nil {
		return false, err
	}
	return IsFilePresent(mmDirPath)
}

func SaveToFile(location string, data []byte) error {
	return os.WriteFile(location, data, 0666)
}
