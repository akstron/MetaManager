package file

import (
	"github/akstron/MetaManager/ds"
	"io/fs"
	"os"
)

/*Get common node information without casting using NodeInformable interface*/
type NodeInformable interface {
	GetAbsPath() string
	GetTags() []string
	AddTag(string)
}
type GeneralNode struct {
	AbsPath string
	Entry   fs.FileInfo
	Tags    []string
}

func NewGeneralNode(absPath string, entry fs.FileInfo) GeneralNode {
	return GeneralNode{
		AbsPath: absPath,
		Entry:   entry,
	}
}

func (gn *GeneralNode) GetAbsPath() string {
	return gn.AbsPath
}

func (gn *GeneralNode) GetTags() []string {
	return gn.Tags
}

func (gn *GeneralNode) AddTag(tag string) {
	for _, lsTag := range gn.Tags {
		if lsTag == tag {
			return
		}
	}

	gn.Tags = append(gn.Tags, tag)
}

type SerializableNode interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}

type FileNode struct {
	GeneralNode
}

func (fn *FileNode) GetInfoProvider() NodeInformable {
	return fn
}

func (fn *FileNode) Name() string {
	return "FILE"
}

type DirNode struct {
	GeneralNode
}

func (dn *DirNode) GetInfoProvider() NodeInformable {
	return dn
}

func (dn *DirNode) Name() string {
	return "DIR"
}

func CreateTreeNodeFromPathAndType(path string, isDir bool) (*ds.TreeNode, error) {
	var info ds.TreeNodeInformable
	if isDir {
		info = &DirNode{
			GeneralNode: GeneralNode{
				AbsPath: path,
			},
		}
	} else {
		info = &FileNode{
			GeneralNode: GeneralNode{
				AbsPath: path,
			},
		}
	}

	return &ds.TreeNode{
		Info: info,
	}, nil
}

func CreateTreeNodeFromPath(path string) (*ds.TreeNode, error) {
	entry, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	return CreateTreeNodeFromPathAndType(path, entry.IsDir())
}
