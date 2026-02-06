package filesys

import (
	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/cmderror"
	"github.com/heroku/self/MetaManager/internal/config"
	"github.com/heroku/self/MetaManager/internal/file"
	"github.com/heroku/self/MetaManager/internal/utils"
	"os"
	"path/filepath"
)

type ScannableCxt map[string]any

type ScannableNode interface {
	// EvalNode would be called at the beginning of every scan operation
	EvalNode(ScannableCxt) error
	ConstructTreeNode() (*ds.TreeNode, error)
	GetChildren() ([]ScannableNode, []ScannableCxt, error)
}

// File System Scanner
type FSScannableNode struct {
	strAbsPath string
	cTreeNode  *ds.TreeNode
	children   []ScannableNode
}

func NewFSScannableNode(absPath string) *FSScannableNode {
	return &FSScannableNode{
		strAbsPath: absPath,
	}
}

func (f *FSScannableNode) GetChildren() ([]ScannableNode, []ScannableCxt, error) {
	scCxt := make([]ScannableCxt, len(f.children))
	return f.children, scCxt, nil
}

func (f *FSScannableNode) ConstructTreeNode() (*ds.TreeNode, error) {
	return f.cTreeNode, nil
}

func (f *FSScannableNode) EvalNode(cxt ScannableCxt) error {
	f.children = []ScannableNode{}

	node, err := os.Stat(f.strAbsPath)
	if err != nil {
		return err
	}

	if node.IsDir() {
		dirNode := &file.DirNode{
			GeneralNode: file.NewGeneralNode(f.strAbsPath, node),
		}
		f.cTreeNode = ds.NewTreeNode(dirNode)

		entries, err := os.ReadDir(f.strAbsPath)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			absEntryPath, err := filepath.Abs(f.strAbsPath + "/" + entry.Name())
			if err != nil {
				return err
			}
			nextNode := &FSScannableNode{
				strAbsPath: absEntryPath,
			}
			f.children = append(f.children, nextNode)
		}
	} else {
		fileNode := &file.FileNode{
			GeneralNode: file.NewGeneralNode(f.strAbsPath, node),
		}
		f.cTreeNode = ds.NewTreeNode(fileNode)
	}

	return nil
}

type GDriveScannableNode struct {
}

func ScanDirectoryV2(dirPath string) (*ds.TreeNode, error) {
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

	scNode := NewFSScannableNode(dirPathAbs)
	scCxt := make(map[string]any)
	return scanDirV2(scNode, scCxt)
}

func scanDirV2(sc ScannableNode, scCxt ScannableCxt) (*ds.TreeNode, error) {
	err := sc.EvalNode(scCxt)
	if err != nil {
		return nil, err
	}

	curTreeNode, err := sc.ConstructTreeNode()
	if err != nil {
		return nil, err
	}

	chs, chCxts, err := sc.GetChildren()
	if err != nil {
		return nil, err
	}

	for i := range chs {
		childTreeNode, err := scanDirV2(chs[i], chCxts[i])
		if err != nil {
			return nil, err
		}
		curTreeNode.AddChild(childTreeNode)
	}

	return curTreeNode, nil
}

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
