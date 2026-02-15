// Package tree provides tree read/write backends (e.g. file-based) for MetaManager.
package tree

import "github.com/heroku/self/MetaManager/internal/ds"

// Fine-grained interface so that we can optimally read trees according
// to the storage system
type TreeReader interface {
	Read() (*ds.TreeNode, error)
}

// Fine-grained interface so that we can optimally write trees according
// to the storage system
type TreeWriter interface {
	Write(*ds.TreeNode) error
}

type TreeRW interface {
	TreeReader
	TreeWriter
}

func WriteTree(tw TreeWriter, node *ds.TreeNode) error {
	return tw.Write(node)
}

func ReadTree(tr TreeReader) (*ds.TreeNode, error) {
	return tr.Read()
}

type RWFactory interface {
	GetTreeRW() (TreeRW, error)
}

/*
This will be changed based on certain flags -> Currently not implemented
*/
func GetTreeRW(factory RWFactory) (TreeRW, error) {
	return factory.GetTreeRW()
}
