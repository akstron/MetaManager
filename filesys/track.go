package filesys

import (
	"github/akstron/MetaManager/ds"
	"github/akstron/MetaManager/pkg/cmderror"
	"github/akstron/MetaManager/pkg/file"
)

/*
Builds a node if it is an absolute location
Builds a subtree if it is an absolute locaiton with "/*"
*/
func Track(path string) (*ds.TreeNode, error) {
	if len(path) == 0 {
		return nil, &cmderror.InvalidPath{}
	}

	if path[len(path)-1] == '*' {
		return trackRec(path)
	}

	return trackAbs(path)
}

func trackAbs(path string) (*ds.TreeNode, error) {
	return file.CreateTreeNodeFromPath(path)
}

func trackRec(path string) (*ds.TreeNode, error) {
	rootDirPath := path[0 : len(path)-1]
	return ScanDirectory(rootDirPath)
}
