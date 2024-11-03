package file

import (
	"io/fs"
)

/*Get common node information without casting using NodeInformable interface*/
type NodeInformable interface {
	GetAbsPath() string
	GetTags() []string
	AddTag(string)
}
type GeneralNode struct {
	absPath string
	entry   fs.FileInfo
	Tags    []string
}

func NewGeneralNode(absPath string, entry fs.FileInfo) GeneralNode {
	return GeneralNode{
		absPath: absPath,
		entry:   entry,
	}
}

func (gn *GeneralNode) GetAbsPath() string {
	return gn.absPath
}

func (gn *GeneralNode) GetTags() []string {
	return gn.Tags
}

func (gn *GeneralNode) AddTag(tag string) {
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

func (fn *FileNode) Scan(ignorable ScanIgnorable) error {
	// Since, this is not scanning anything, no requirement for check ignorable
	return nil
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
