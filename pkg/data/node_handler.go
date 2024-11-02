package data

type ScanningHandler interface {
	HandleFile(*TreeNode, *FileNode) error
	HandleDir(*TreeNode, *DirNode) error
}

// type ScanningHandler interface {
// 	Handle(*TreeNode, any) error
// }

type ScanHandler struct {
	ig ScanIgnorable
}

func NewScanHandler(ig ScanIgnorable) *ScanHandler {
	return &ScanHandler{
		ig: ig,
	}
}

func (sh *ScanHandler) HandleFile(parentNode *TreeNode, curNode *FileNode) error {
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

	parentNode.children = append(parentNode.children, &TreeNode{info: curNode})
	return nil
}

func (sh *ScanHandler) HandleDir(parentNode *TreeNode, curNode *DirNode) error {
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

	parentNode.children = append(parentNode.children, &TreeNode{info: curNode})
	return nil
}
