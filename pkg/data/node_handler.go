package data

type ScanningHandler interface {
	HandleFile(*DirNode, *FileNode) error
	HandleDir(*DirNode, *DirNode) error
}

type ScanHandler struct {
	ig ScanIgnorable
}

func (sh *ScanHandler) HandleFile(parentNode *DirNode, curNode *FileNode) error {
	/*
		Ignorer checks if this path should be ignored
	*/
	shouldIgnore, err := sh.ig.ShouldIgnore(curNode.absPath)
	if err != nil {
		return err
	}

	if shouldIgnore {
		return nil
	}

	parentNode.FileChildren = append(parentNode.FileChildren, curNode)
	return nil
}

func (sh *ScanHandler) HandleDir(parentNode, curNode *DirNode) error {
	/*
		Ignorer checks if this path should be ignored
	*/
	shouldIgnore, err := sh.ig.ShouldIgnore(curNode.absPath)
	if err != nil {
		return err
	}

	if shouldIgnore {
		return nil
	}

	parentNode.DirChildren = append(parentNode.DirChildren, curNode)
	return nil
}

// TODO: Use it later
type Iterable interface {
	Next() error
	HasNext() error
}
