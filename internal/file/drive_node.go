package file

import (
	"strings"

	"github.com/heroku/self/MetaManager/internal/ds"
)

// GDrivePathPrefix is the virtual path prefix for Google Drive nodes (e.g. "gdrive:/Folder/file").
const GDrivePathPrefix = "gdrive:/"
const GDrivePathRoot = "gdrive:"

// IsGDrivePath returns true if path is a Drive virtual path.
func IsGDrivePath(path string) bool {
	return strings.HasPrefix(path, GDrivePathPrefix)
}

// NewDriveDirNode creates a tree node for a Drive folder.
func NewDriveDirNode(virtualPath, driveId string) *ds.TreeNode {
	return ds.NewTreeNode(&FileNode{
		GeneralNode: GeneralNode{AbsPath: virtualPath},
		DriveId:     driveId,
	})
}

// NewDriveFileNode creates a tree node for a Drive file.
func NewDriveFileNode(virtualPath, driveId string) *ds.TreeNode {
	return ds.NewTreeNode(&FileNode{
		GeneralNode: GeneralNode{AbsPath: virtualPath},
		DriveId:     driveId,
	})
}

func CreateTreeNodeFromPath(path string) (*ds.TreeNode, error) {
	// Normal file node
	return &ds.TreeNode{
		Info: &FileNode{
			GeneralNode: GeneralNode{AbsPath: path},
		},
	}, nil
}
