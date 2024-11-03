package file

import "github/akstron/MetaManager/ds"

type ScanningHandler interface {
	// HandleFile(*ds.TreeNode, *FileNode) error
	// HandleDir(*ds.TreeNode, *DirNode) error
	Handle(*ds.TreeNode, *ds.TreeNode) error
}

// type ScanningHandler interface {
// 	Handle(*ds.TreeNode, any) error
// }

type ScanHandler struct {
	ig ScanIgnorable
}

func NewScanHandler(ig ScanIgnorable) *ScanHandler {
	return &ScanHandler{
		ig: ig,
	}
}

func (sh *ScanHandler) Handle(parentNode *ds.TreeNode, curNode *ds.TreeNode) error {
	parentNode.AddChild(curNode)

	return nil
}

// func (sh *ScanHandler) HandleFile(parentNode *ds.TreeNode, curNode *FileNode) error {
// 	/*
// 		Ignorer checks if this path should be ignored
// 	*/
// 	shouldIgnore, err := sh.ig.ShouldIgnore(curNode.absPath)
// 	if err != nil {
// 		return err
// 	}

// 	if shouldIgnore {
// 		return nil
// 	}

// 	parentNode.AddChild(ds.NewTreeNode(curNode))

// 	return nil
// }

// func (sh *ScanHandler) HandleDir(parentNode *ds.TreeNode, curNode *DirNode) error {
// 	/*
// 		Ignorer checks if this path should be ignored
// 	*/
// 	shouldIgnore, err := sh.ig.ShouldIgnore(curNode.absPath)
// 	if err != nil {
// 		return err
// 	}

// 	if shouldIgnore {
// 		return nil
// 	}

// 	parentNode.AddChild(ds.NewTreeNode(curNode))
// 	return nil
// }
