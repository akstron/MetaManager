package filesys

import (
	"github/akstron/MetaManager/ds"
	"github/akstron/MetaManager/pkg/cmderror"
	"github/akstron/MetaManager/pkg/config"
	"github/akstron/MetaManager/pkg/file"
	"github/akstron/MetaManager/pkg/utils"
	"os"
	"path/filepath"
)

func ScanDirectory(dirPath string) (*ds.TreeNode, error) {
	present, err := utils.IsFilePresent(dirPath)
	if err != nil {
		return nil, err
	}

	if !present {
		return nil, &cmderror.InvalidPath{}
	}

	dirPathAbs, err := filepath.Abs(dirPath)
	if err != nil {
		return nil, err
	}

	root, err := os.Stat(dirPathAbs)
	if err != nil {
		return nil, err
	}

	topDir := &file.DirNode{GeneralNode: file.NewGeneralNode(dirPathAbs, root)}
	// treeNode := data.NewTreeNode(topDir)
	// trMg := data.NewTreeManager(&treeNode)

	igMg, err := config.NewIgnoreManager()
	if err != nil {
		return nil, err
	}

	igMg.Load()
	ignorer := file.NewNodeAbsPathIgnorer(igMg)

	/*
		TODO: Probably better to pass pointer
		We will use go routines to accelerate dfs
	*/
	handler := file.NewScanHandler(ignorer)

	topTreeNode, err := ScanDir(topDir, handler)
	if err != nil {
		return nil, err
	}

	return topTreeNode, nil
}

func ScanFile(fn *file.FileNode, handler file.ScanningHandler) (*ds.TreeNode, error) {
	return &ds.TreeNode{
		Info: fn,
	}, nil
}

func ScanDir(fn *file.DirNode, handler file.ScanningHandler) (*ds.TreeNode, error) {
	curTreeNode := &ds.TreeNode{
		Info: fn,
	}

	entries, err := os.ReadDir(fn.GetAbsPath())
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		// var curNode Node
		fileEntry, err := entry.Info()
		if err != nil {
			return nil, err
		}

		absEntryPath, err := filepath.Abs(fn.GetAbsPath() + "/" + entry.Name())
		if err != nil {
			return nil, err
		}

		// TODO: Implement some common function (probably)
		if entry.IsDir() {
			dirNode := &file.DirNode{
				GeneralNode: file.NewGeneralNode(absEntryPath, fileEntry),
			}

			if childNode, err := ScanDir(dirNode, handler); err != nil {
				return nil, err
			} else {
				err = handler.Handle(curTreeNode, childNode)
				if err != nil {
					return nil, err
				}
			}
		} else {
			fileNode := &file.FileNode{
				GeneralNode: file.NewGeneralNode(absEntryPath, fileEntry),
			}

			if childNode, err := ScanFile(fileNode, handler); err != nil {
				return nil, err
			} else {
				err = handler.Handle(curTreeNode, childNode)
				if err != nil {
					return nil, err
				}
			}
		}
		// TODO: Convert this to BFS
	}

	return curTreeNode, nil
}
