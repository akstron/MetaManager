package utils

import "github.com/heroku/self/MetaManager/internal/cmderror"

// CommonInitChecks returns true if the given context's .mm directory exists.
func CommonInitChecks(contextName string) (bool, error) {
	isInitialized, err := IsRootInitialized(contextName)
	if err != nil {
		return false, err
	}
	return isInitialized, nil
}

// CommonAlreadyInitializedChecks returns the .mm directory path for the given context, or error if not set/uninitialized.
func CommonAlreadyInitializedChecks(contextName string) (string, error) {
	if contextName == "" {
		return "", &cmderror.UninitializedRoot{}
	}
	isPresent, dirPath, err := FindMMDirPath(contextName)
	if err != nil {
		return "", err
	}
	if !isPresent {
		return "", &cmderror.UninitializedRoot{}
	}
	return dirPath, nil
}
