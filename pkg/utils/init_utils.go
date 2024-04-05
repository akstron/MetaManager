package utils

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
