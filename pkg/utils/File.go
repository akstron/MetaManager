package utils

import "os"

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
