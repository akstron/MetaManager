package filesys

import (
	"github/akstron/MetaManager/pkg/cmderror"
	"github/akstron/MetaManager/pkg/config"
	"github/akstron/MetaManager/pkg/data"
	"github/akstron/MetaManager/pkg/utils"
	"os"
	"path/filepath"
)

func ScanDirectory(dirPath string) (*data.TreeNode, error) {
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

	topDir := &data.DirNode{GeneralNode: data.NewGeneralNode(dirPathAbs, root)}
	// treeNode := data.NewTreeNode(topDir)
	// trMg := data.NewTreeManager(&treeNode)

	igMg, err := config.NewIgnoreManager()
	if err != nil {
		return nil, err
	}

	igMg.Load()
	ignorer := data.NewNodeAbsPathIgnorer(igMg)

	/*
		TODO: Probably better to pass pointer
		We will use go routines to accelerate dfs
	*/
	handler := data.NewScanHandler(ignorer)

	topTreeNode, err := data.ScanDir(topDir, handler)
	if err != nil {
		return nil, err
	}

	return topTreeNode, nil
}
