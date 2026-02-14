package filesys

import (
	"github.com/heroku/self/MetaManager/internal/ds"
)

type ScanningHandler interface {
	Handle(*ds.TreeNode, *ds.TreeNode) error
}

type ScanHandler struct{}

func NewScanHandler() *ScanHandler {
	return &ScanHandler{}
}

func (sh *ScanHandler) Handle(parentNode *ds.TreeNode, curNode *ds.TreeNode) error {
	parentNode.AddChild(curNode)
	return nil
}
