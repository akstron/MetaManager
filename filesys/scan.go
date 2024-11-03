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

	igMg, err := config.NewIgnoreManager()
	if err != nil {
		return nil, err
	}

	igMg.Load()
	ignorer := NewNodeAbsPathIgnorer(igMg)

	/*
		TODO: Probably better to pass pointer
		We will use go routines to accelerate dfs
	*/
	handler := NewScanHandler(ignorer)

	topTreeNode, err := scanDir(topDir, handler)
	if err != nil {
		return nil, err
	}

	return topTreeNode, nil
}

func scanFile(fn *file.FileNode, handler ScanningHandler) (*ds.TreeNode, error) {
	return &ds.TreeNode{
		Info: fn,
	}, nil
}

func scanDir(fn *file.DirNode, handler ScanningHandler) (*ds.TreeNode, error) {
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

			if childNode, err := scanDir(dirNode, handler); err != nil {
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

			if childNode, err := scanFile(fileNode, handler); err != nil {
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
