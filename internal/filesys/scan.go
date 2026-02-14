package filesys

import (
	"os"
	"path/filepath"

	"github.com/heroku/self/MetaManager/internal/cmderror"
	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/file"
	"github.com/heroku/self/MetaManager/internal/utils"
)

type Scanner interface {
	Scan(path string) (*ds.TreeNode, error)
}

func NewUnixFileSystemScanner() *UnixFileSystemScanner {
	return &UnixFileSystemScanner{}
}

type UnixFileSystemScanner struct {
}

type ScannableCxt map[string]any

type ScannableNode interface {
	// EvalNode would be called at the beginning of every scan operation
	EvalNode(ScannableCxt) error
	ConstructTreeNode() (*ds.TreeNode, error)
	GetChildren() ([]ScannableNode, []ScannableCxt, error)
}

// TODO: Implement this
type GDriveScannableNode struct {
	strAbsPath string
	cTreeNode  *ds.TreeNode
	children   []ScannableNode
}

func NewGDriveScannableNode(absPath string) *GDriveScannableNode {
	return &GDriveScannableNode{
		strAbsPath: absPath,
	}
}

func (g *GDriveScannableNode) EvalNode(cxt ScannableCxt) error {
	return nil
}

func (g *GDriveScannableNode) ConstructTreeNode() (*ds.TreeNode, error) {
	return g.cTreeNode, nil
}

func (g *GDriveScannableNode) GetChildren() ([]ScannableNode, []ScannableCxt, error) {
	return g.children, nil, nil
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
		nodeInfo := &file.FileNode{
			GeneralNode: file.NewGeneralNode(f.strAbsPath, node),
		}
		f.cTreeNode = ds.NewTreeNode(nodeInfo)

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
		nodeInfo := &file.FileNode{
			GeneralNode: file.NewGeneralNode(f.strAbsPath, node),
		}
		f.cTreeNode = ds.NewTreeNode(nodeInfo)
	}

	return nil
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
