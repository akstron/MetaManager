package filesys

import (
	"github.com/heroku/self/MetaManager/internal/cmderror"
	"github.com/heroku/self/MetaManager/internal/ds"
)

// Scanner is the interface for scanning directories and returning tree nodes.
type Scanner interface {
	Scan(path string) (*ds.TreeNode, error)
}

// CreateScannerFromContextType creates a scanner based on the context type.
func CreateScannerFromContextType(contextType string) (Scanner, error) {
	switch contextType {
	case "local":
		return NewUnixFileSystemScanner(), nil
	case "gdrive":
		return createGDriveScanner()
	}
	return nil, &cmderror.InvalidOperation{}
}

// ScannableCxt is the context passed to scannable nodes during evaluation.
type ScannableCxt map[string]any

// ScannableNode is the interface for nodes that can be scanned.
type ScannableNode interface {
	// EvalNode would be called at the beginning of every scan operation
	EvalNode(ScannableCxt) error
	ConstructTreeNode() (*ds.TreeNode, error)
	GetChildren() ([]ScannableNode, []ScannableCxt, error)
}
