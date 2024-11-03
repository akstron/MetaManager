package data

import (
	"fmt"
	"github/akstron/MetaManager/ds"
	"github/akstron/MetaManager/pkg/cmderror"
	"github/akstron/MetaManager/pkg/file"
	"path/filepath"
	"slices"
)

type DirTreeManager struct {
	*ds.TreeManager
}

func NewDirTreeManager(trMg *ds.TreeManager) *DirTreeManager {
	return &DirTreeManager{
		TreeManager: trMg,
	}
}

func (mg *DirTreeManager) MergeNodeWithPath(path string) error {
	treeNode, err := file.CreateTreeNodeFromPath(path)
	if err != nil {
		return err
	}

	return mg.MergeNode(treeNode)
}

func (mg *DirTreeManager) MergeNode(treeNode *ds.TreeNode) error {
	if treeNode == nil {
		return &cmderror.InvalidOperation{}
	}

	if mg.TreeManager == nil {
		mg.TreeManager = ds.NewTreeManager(treeNode)
		return nil
	}

	if mg.TreeManager.Root == nil {
		mg.TreeManager.Root = treeNode
		return nil
	}

	fir, ok := mg.Root.Info.(file.NodeInformable)
	if !ok {
		return &cmderror.Unexpected{}
	}

	iter := ds.NewTreeIterator(ds.NewTreeManager(treeNode))
	for iter.HasNext() {
		curNode, err := iter.Next()
		got := curNode.Info
		if err != nil {
			return err
		}

		sec, ok := got.(file.NodeInformable)
		if !ok {
			return &cmderror.Unexpected{}
		}

		// Absolute path of fir should be a prefix of absolute path
		midPaths := make([]string, 0)
		firPath := fir.GetAbsPath()
		secPath := sec.GetAbsPath()

		for firPath != secPath && len(secPath) > 0 {
			midPaths = append(midPaths, secPath)
			secPath = filepath.Join(secPath, "..")
		}

		slices.Reverse(midPaths)

		err = mg.createPathNodes(midPaths)
		if err != nil {
			return err
		}
	}

	return nil
}

func (mg *DirTreeManager) createPathNodes(paths []string) error {
	return mg.createPathNodesInternal(mg.Root, paths, 0)
}

func (mg *DirTreeManager) createPathNodesInternal(curNode *ds.TreeNode, paths []string, index int) error {
	if curNode == nil {
		return &cmderror.InvalidOperation{}
	}

	var err error
	if index >= len(paths) {
		return nil
	}

	reqPath := paths[index]

	var nextNode *ds.TreeNode
	for _, child := range curNode.Children {
		info, ok := child.Info.(file.NodeInformable)
		if !ok {
			return &cmderror.Unexpected{}
		}
		if info.GetAbsPath() == reqPath {
			nextNode = child
			break
		}
	}

	if nextNode == nil {
		nextNode, err = file.CreateTreeNodeFromPath(reqPath)
		if err != nil {
			return err
		}
		curNode.Children = append(curNode.Children, nextNode)
	}

	return mg.createPathNodesInternal(nextNode, paths, index+1)
}

func (mg *DirTreeManager) FindNodeByAbsPath(path string) (file.NodeInformable, error) {
	ti := ds.NewTreeIterator(mg.TreeManager)
	return mg.findNodeByAbsPathInternal(ti, path)
}

func (mg *DirTreeManager) findNodeByAbsPathInternal(it ds.TreeIterator, path string) (file.NodeInformable, error) {
	for it.HasNext() {
		curNode, err := it.Next()
		got := curNode.Info
		if err != nil {
			return nil, err
		}

		if got.(file.NodeInformable).GetAbsPath() == path {
			return got.(file.NodeInformable), nil
		}
	}

	return nil, fmt.Errorf("node not found")
}
