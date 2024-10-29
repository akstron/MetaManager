package utils

import "github/akstron/MetaManager/pkg/cmderror"

// Did not used it as it may hide errors
func CommonInitChecks() (bool, error) {
	isInitialized, err := IsRootInitialized()
	if err != nil {
		return false, err
	}

	if !isInitialized {
		return false, nil
	}

	return true, nil
}

func CommonAlreadyInitializedChecks() (string, error) {
	isPresent, dirPath, err := FindMMDirPath()
	if err != nil {
		return "", err
	}

	if !isPresent {
		return "", &cmderror.UninitializedRoot{}
	}

	return dirPath, nil
}
