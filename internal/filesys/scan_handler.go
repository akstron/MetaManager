package filesys

import (
	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/config"
)

/*
Why this interface?
Later it would be helpful to ignore on the basis of prefix/suffix etc
by having different implementations
Or maybe we would like to ignore based on REGEX
*/
type ScanIgnorable interface {
	ShouldIgnore(string) (bool, error)
}

type NodeAbsPathIgnorer struct {
	igMg *config.IgnoreManager
}

func NewNodeAbsPathIgnorer(igMg *config.IgnoreManager) *NodeAbsPathIgnorer {
	return &NodeAbsPathIgnorer{
		igMg: igMg,
	}
}

func (ig *NodeAbsPathIgnorer) ShouldIgnore(ignorePath string) (bool, error) {
	/*
		igMg can have a GetData which returns constant data for iteration
		but lets see if this should be done
	*/
	for _, value := range ig.igMg.Data.Paths {
		if value == ignorePath {
			return true, nil
		}
	}
	return false, nil
}

type ScanningHandler interface {
	Handle(*ds.TreeNode, *ds.TreeNode) error
}

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
