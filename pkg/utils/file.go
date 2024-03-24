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

func GetAbsRootPath() (string, error) {
	dirPath := "./.mm"
	return filepath.Abs(dirPath)
}
