package filesys

import (
	"os"
	"path/filepath"

	"github.com/heroku/self/MetaManager/internal/cmderror"
	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/file"
	"github.com/heroku/self/MetaManager/internal/utils"
)

// UnixFileSystemScanner scans local file system directories.
type UnixFileSystemScanner struct {
}

// NewUnixFileSystemScanner creates a new UnixFileSystemScanner.
func NewUnixFileSystemScanner() *UnixFileSystemScanner {
	return &UnixFileSystemScanner{}
}

// Scan scans a directory path and returns a tree node.
func (u *UnixFileSystemScanner) Scan(path string) (*ds.TreeNode, error) {
	return ScanDirectoryV2(path)
}

// Make sure that UnixFileSystemScanner implements Scanner
var _ Scanner = (*UnixFileSystemScanner)(nil)

// FSScannableNode is a scannable node for the file system.
type FSScannableNode struct {
	strAbsPath string
	cTreeNode  *ds.TreeNode
	children   []ScannableNode
}

// NewFSScannableNode creates a new file system scannable node.
func NewFSScannableNode(absPath string) *FSScannableNode {
	return &FSScannableNode{
		strAbsPath: absPath,
	}
}

// GetChildren returns the children nodes and their contexts.
func (f *FSScannableNode) GetChildren() ([]ScannableNode, []ScannableCxt, error) {
	scCxt := make([]ScannableCxt, len(f.children))
	return f.children, scCxt, nil
}

// ConstructTreeNode returns the constructed tree node.
func (f *FSScannableNode) ConstructTreeNode() (*ds.TreeNode, error) {
	return f.cTreeNode, nil
}

// EvalNode evaluates the node by reading file system information.
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

// ScanDirectoryV2 scans a directory and returns a tree node representation.
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

// scanDirV2 recursively scans a scannable node and its children.
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
