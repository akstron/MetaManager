package data

import (
	"fmt"
	"github/akstron/MetaManager/ds"
	"github/akstron/MetaManager/pkg/file"
	"os"
)

type DirTreeManager struct {
	*ds.TreeManager
}

func CreateTreeNodeFromPath(path string) (*ds.TreeNode, error) {
	entry, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	var info ds.TreeNodeInformable
	if entry.IsDir() {
		info = &file.DirNode{
			GeneralNode: file.GeneralNode{
				AbsPath: path,
				Entry:   entry,
			},
		}
	} else {
		info = &file.FileNode{
			GeneralNode: file.GeneralNode{
				AbsPath: path,
				Entry:   entry,
			},
		}
	}

	return &ds.TreeNode{
		Info: info,
	}, nil
}

func (mg *DirTreeManager) MergeNodeWithPath(path string) error {
	treeNode, err := CreateTreeNodeFromPath(path)
	if err != nil {
		return err
	}

	return mg.MergeNode(treeNode)
}

func (mg *DirTreeManager) MergeNode(treeNode *ds.TreeNode) error {
	if mg.TreeManager == nil {
		mg.TreeManager = ds.NewTreeManager(treeNode)
		return nil
	}
	return fmt.Errorf("not implemented")
}

func (mg *DirTreeManager) FindNodeByAbsPath(path string) (file.NodeInformable, error) {
	ti := ds.NewTreeIterator(mg.TreeManager)
	return mg.findNodeByAbsPathInternal(ti, path)
}

func (mg *DirTreeManager) findNodeByAbsPathInternal(it ds.TreeIterator, path string) (file.NodeInformable, error) {
	for it.HasNext() {
		got, err := it.Next()
		if err != nil {
			return nil, err
		}

		if got.(file.NodeInformable).GetAbsPath() == path {
			return got.(file.NodeInformable), nil
		}
	}

	return nil, nil
}
