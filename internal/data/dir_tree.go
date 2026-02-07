package data

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"slices"

	"github.com/heroku/self/MetaManager/internal/cmderror"
	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/file"
	"github.com/sirupsen/logrus"
)

type DirTreeManager struct {
	*ds.TreeManager
}

func NewDirTreeManager(trMg *ds.TreeManager) *DirTreeManager {
	return &DirTreeManager{
		TreeManager: trMg,
	}
}

// We probably don't want to update anything on the original extracted tree nodes
// So, we create a copy of each tree node and buildTree out of it
func BuildCopyTree(rootPath string, treeNodes []*ds.TreeNode) (*DirTreeManager, error) {
	copyTreeNodes := []*ds.TreeNode{}
	rootNode, err := file.CreateTreeNodeFromPath(rootPath)
	if err != nil {
		return nil, err
	}

	for _, node := range treeNodes {
		copyNode := *node
		copyTreeNodes = append(copyTreeNodes, &copyNode)
	}

	return buildTree(rootNode, copyTreeNodes)
}

// builds a tree out of the nodes
func buildTree(rootNode *ds.TreeNode, nodes []*ds.TreeNode) (*DirTreeManager, error) {
	drMg := NewDirTreeManager(ds.NewTreeManager(rootNode))
	for _, node := range nodes {
		err := drMg.MergeNode(node)
		if err != nil {
			return nil, err
		}
	}
	return drMg, nil
}

func (mg *DirTreeManager) FindTreeNodesByRegex(expression string) ([]*ds.TreeNode, error) {
	it := ds.NewTreeIterator(mg.TreeManager)
	// iterate over all the node and find node with given regex
	return mg.findTreeNodesByRegexInternal(expression, it)
}

func (mg *DirTreeManager) findTreeNodesByRegexInternal(pattern string, it ds.TreeIterable) ([]*ds.TreeNode, error) {
	nodesFound := []*ds.TreeNode{}

	for it.HasNext() {
		curNode, err := it.Next()
		got := curNode.Info
		if err != nil {
			return nil, err
		}

		nodeInfo, ok := got.(file.NodeInformable)
		if !ok {
			return nil, errors.New("info not convertiable to NodeInformable")
		}

		match, err := regexp.MatchString(pattern, nodeInfo.GetAbsPath())
		if err != nil {
			return nil, err
		}

		if match {
			nodesFound = append(nodesFound, curNode)
		}
	}

	return nodesFound, nil
}

func (mg *DirTreeManager) SplitChildrenFromPath(path string) error {
	curTreeNode, err := mg.FindTreeNodeByAbsPath(path)
	if err != nil {
		return err
	}
	curTreeNode.Children = []*ds.TreeNode{}
	return nil
}

func (mg *DirTreeManager) SplitNodeWithPath(path string) error {
	if mg.Root.Info.(file.NodeInformable).GetAbsPath() == path {
		mg.Root = nil
		return nil
	}

	parentPath := filepath.Join(path, "..")
	parentTreeNode, err := mg.FindTreeNodeByAbsPath(parentPath)
	if err != nil {
		return err
	}

	newChildren := []*ds.TreeNode{}
	for _, child := range parentTreeNode.Children {
		info, ok := child.Info.(file.NodeInformable)
		if !ok {
			return &cmderror.Unexpected{}
		}
		if info.GetAbsPath() != path {
			newChildren = append(newChildren, child)
		}
	}
	parentTreeNode.Children = newChildren

	return nil
}

// Given a path of any form /a/b/c -> this will create tree out of this path and attach it to the tree managed by mg
func (mg *DirTreeManager) MergeNodeWithPath(path string) error {
	treeNode, err := file.CreateTreeNodeFromPath(path)
	if err != nil {
		return err
	}

	return mg.MergeNode(treeNode)
}

// isPathPrefixOrEqual returns true when root is equal to path or path is under root (root is a path prefix of path).
func isPathPrefixOrEqual(root, path string) bool {
	root = filepath.Clean(root)
	path = filepath.Clean(path)
	if root == path {
		return true
	}
	// Ensure we don't match /foo against /foobar: path must start with root followed by separator.
	return len(path) > len(root) && (path[len(root)] == '/' || path[len(root)] == filepath.Separator) && path[:len(root)] == root
}

func (mg *DirTreeManager) MergeNode(treeNode *ds.TreeNode) error {
	logrus.Debugf("[merge] MergeNode start treeNode=%v", treeNode != nil)
	if treeNode == nil {
		return &cmderror.InvalidOperation{}
	}

	if mg.TreeManager == nil {
		logrus.Debugf("[merge] TreeManager nil, setting root to incoming tree")
		mg.TreeManager = ds.NewTreeManager(treeNode)
		return nil
	}

	if mg.TreeManager.Root == nil {
		logrus.Debugf("[merge] Root nil, setting root to incoming tree")
		mg.TreeManager.Root = treeNode
		return nil
	}

	fir, ok := mg.Root.Info.(file.NodeInformable)
	if !ok {
		return &cmderror.Unexpected{}
	}
	rootPath := fir.GetAbsPath()
	logrus.Debugf("[merge] current root path=%q", rootPath)

	iter := ds.NewTreeIterator(ds.NewTreeManager(treeNode))
	nodeCount := 0
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
		secPathOrig := sec.GetAbsPath()
		secPath := secPathOrig
		nodeCount++

		// Absolute path of root (firPath) must be a prefix of the node path (secPath), or equal.
		if !isPathPrefixOrEqual(rootPath, secPathOrig) {
			return fmt.Errorf("merge: root path %q is not a prefix of node path %q", rootPath, secPathOrig)
		}

		midPaths := make([]string, 0)
		firPath := rootPath

		for firPath != secPath && len(secPath) > 0 {
			midPaths = append(midPaths, secPath)
			secPath = filepath.Join(secPath, "..")
		}

		slices.Reverse(midPaths)
		logrus.Debugf("[merge] node %d path=%q midPaths=%v", nodeCount, secPathOrig, midPaths)

		err = mg.createPathNodes(midPaths)
		if err != nil {
			logrus.Debugf("[merge] createPathNodes error: %v", err)
			return err
		}
	}

	logrus.Debugf("[merge] MergeNode done, merged %d nodes", nodeCount)
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
	curPath := ""
	if info, ok := curNode.Info.(file.NodeInformable); ok {
		curPath = info.GetAbsPath()
	}

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
		logrus.Debugf("[merge] createPathNodes index=%d curPath=%q creating node for %q", index, curPath, reqPath)
		nextNode, err = file.CreateTreeNodeFromPath(reqPath)
		if err != nil {
			logrus.Debugf("[merge] CreateTreeNodeFromPath %q error: %v", reqPath, err)
			return err
		}
		curNode.Children = append(curNode.Children, nextNode)
	}

	return mg.createPathNodesInternal(nextNode, paths, index+1)
}

func (mg *DirTreeManager) FindFileNodeById(id string) (file.NodeInformable, error) {
	node, err := mg.FindTreeNodeById(id)
	if err != nil {
		return nil, err
	}

	nodeInfo, ok := node.Info.(file.NodeInformable)
	if !ok {
		return nil, &cmderror.Unexpected{}
	}
	return nodeInfo, nil
}

func (mg *DirTreeManager) FindTreeNodeById(id string) (*ds.TreeNode, error) {
	ti := ds.NewTreeIterator(mg.TreeManager)
	return mg.findTreeNodeByIdInternal(id, ti)
}

func (mg *DirTreeManager) findTreeNodeByIdInternal(id string, it ds.TreeIterable) (*ds.TreeNode, error) {
	for it.HasNext() {
		curNode, err := it.Next()
		got := curNode.Info
		if err != nil {
			return nil, err
		}

		if nodeInfo, ok := got.(file.NodeInformable); ok {
			if nodeInfo.GetId() == id {
				return curNode, nil
			}
		} else {
			return nil, errors.New("info not convertiable to NodeInformable")
		}
	}

	return nil, errors.New("node not found")
}

func (mg *DirTreeManager) FindTreeNodeByAbsPath(path string) (*ds.TreeNode, error) {
	ti := ds.NewTreeIterator(mg.TreeManager)
	return mg.findTreeNodeByAbsPathInternal(ti, path)
}

func (mg *DirTreeManager) FindNodeByAbsPath(path string) (file.NodeInformable, error) {
	ti := ds.NewTreeIterator(mg.TreeManager)
	trNode, err := mg.findTreeNodeByAbsPathInternal(ti, path)
	if err != nil {
		return nil, err
	}
	return trNode.Info.(file.NodeInformable), nil
}

func (mg *DirTreeManager) findTreeNodeByAbsPathInternal(it ds.TreeIterable, path string) (*ds.TreeNode, error) {
	for it.HasNext() {
		curNode, err := it.Next()
		got := curNode.Info
		if err != nil {
			return nil, err
		}

		if got.(file.NodeInformable).GetAbsPath() == path {
			return curNode, nil
		}
	}

	return nil, fmt.Errorf("node not found")
}
