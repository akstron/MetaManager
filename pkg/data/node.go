package data

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
)

type GeneralNode struct {
	absPath string
	entry   fs.FileInfo
	Tags    []string
}

/*Get common node information without casting using NodeInformable interface*/
type NodeInformable interface {
	GetAbsPath() string
	GetTags() []string
	AddTag(string)
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

func (fn *FileNode) GetFileChildren() []*FileNode {
	return nil
}

func (fn *FileNode) GetDirChildren() []*DirNode {
	return nil
}

func (fn *FileNode) GetInfoProvider() NodeInformable {
	return &fn.GeneralNode
}

func (fn *FileNode) Scan(ignorable ScanIgnorable) error {
	// Since, this is not scanning anything, no requirement for check ignorable
	return nil
}

type DirNode struct {
	GeneralNode
	DirChildren  []*DirNode
	FileChildren []*FileNode
}

func (dn *DirNode) GetFileChildren() []*FileNode {
	return dn.FileChildren
}

func (dn *DirNode) GetDirChildren() []*DirNode {
	return dn.DirChildren
}

func (dn *DirNode) GetInfoProvider() NodeInformable {
	return &dn.GeneralNode
}

func (fn *DirNode) Scan(handler ScanHandler) error {
	entries, err := os.ReadDir(fn.absPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		// var curNode Node
		fileEntry, err := entry.Info()
		if err != nil {
			return err
		}

		absEntryPath, err := filepath.Abs(fn.absPath + "/" + entry.Name())
		if err != nil {
			return err
		}

		// TODO: Implement some common function (probably)
		if entry.IsDir() {
			dirNode := &DirNode{
				GeneralNode: GeneralNode{
					entry:   fileEntry,
					absPath: absEntryPath,
				},
			}

			err = handler.HandleDir(fn, dirNode)
			if err != nil {
				return err
			}

			if err := dirNode.Scan(handler); err != nil {
				return err
			}
		} else {
			fileNode := &FileNode{
				GeneralNode: GeneralNode{
					entry:   fileEntry,
					absPath: absEntryPath,
				},
			}
			err = handler.HandleFile(fn, fileNode)
			if err != nil {
				return err
			}
		}
		// TODO: Convert this to BFS
	}

	return nil
}

type NodeJSON struct {
	Parent   string
	Children [][]byte
	IsDir    bool
	Tags     []string
}

func (fn *FileNode) MarshalJSON() ([]byte, error) {
	obj := NodeJSON{
		Parent: fn.absPath,
		IsDir:  false,
		Tags:   fn.Tags,
	}
	return json.Marshal(obj)
}

func (fn *FileNode) UnmarshalJSON(data []byte) error {
	var obj NodeJSON
	return json.Unmarshal(data, &obj)
}

func (dn *DirNode) MarshalJSON() ([]byte, error) {
	var childrenSerialized [][]byte
	for _, node := range dn.DirChildren {
		currentChildSerialized, err := json.Marshal(node)
		if err != nil {
			return nil, err
		}
		childrenSerialized = append(childrenSerialized, (currentChildSerialized))
	}

	for _, node := range dn.FileChildren {
		currentChildSerialized, err := json.Marshal(node)
		if err != nil {
			return nil, err
		}
		childrenSerialized = append(childrenSerialized, (currentChildSerialized))
	}

	obj := NodeJSON{
		Parent:   dn.absPath,
		Children: childrenSerialized,
		IsDir:    true,
		Tags:     dn.Tags,
	}

	return json.Marshal(&obj)
}

func (dn *DirNode) UnmarshalJSON(data []byte) error {
	var obj NodeJSON
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return err
	}

	dn.absPath = obj.Parent
	dn.Tags = obj.Tags
	for _, child := range obj.Children {
		var childObj NodeJSON
		err := json.Unmarshal(child, &childObj)
		if err != nil {
			return err
		}

		childAbsPath := childObj.Parent
		var curNode SerializableNode
		if obj.IsDir {
			dirNode := &DirNode{
				GeneralNode: GeneralNode{
					absPath: childAbsPath,
				},
			}
			dn.DirChildren = append(dn.DirChildren, dirNode)
			curNode = dirNode
		} else {
			fileNode := &FileNode{
				GeneralNode: GeneralNode{
					absPath: childAbsPath,
				},
			}
			dn.FileChildren = append(dn.FileChildren, fileNode)
			curNode = fileNode
		}
		err = json.Unmarshal(child, &curNode)
		if err != nil {
			return err
		}
	}

	return nil
}
