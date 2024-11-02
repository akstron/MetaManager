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

func NewGeneralNode(absPath string, entry fs.FileInfo) GeneralNode {
	return GeneralNode{
		absPath: absPath,
		entry:   entry,
	}
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

// func (fn *FileNode) GetFileChildren() []*FileNode {
// 	return nil
// }

// func (fn *FileNode) GetDirChildren() []*DirNode {
// 	return nil
// }

// func (fn *FileNode) GetChildren() []NodeIterable {
// 	return nil
// }

func (fn *FileNode) GetInfoProvider() NodeInformable {
	return fn
}

func (fn *FileNode) Scan(ignorable ScanIgnorable) error {
	// Since, this is not scanning anything, no requirement for check ignorable
	return nil
}

type DirNode struct {
	GeneralNode
	// DirChildren  []*DirNode
	// FileChildren []*FileNode
}

// func (dn *DirNode) GetFileChildren() []*FileNode {
// 	return dn.FileChildren
// }

// func (dn *DirNode) GetDirChildren() []*DirNode {
// 	return dn.DirChildren
// }

// func (dn *DirNode) GetChildren() []NodeIterable {
// 	var result []NodeIterable
// 	for _, node := range dn.GetDirChildren() {
// 		result = append(result, node)
// 	}
// 	for _, node := range dn.GetFileChildren() {
// 		result = append(result, node)
// 	}
// 	return result
// }

func (dn *DirNode) GetInfoProvider() NodeInformable {
	return dn
}

func ScanFile(fn *FileNode, handler ScanningHandler) (*TreeNode, error) {
	return &TreeNode{
		info: fn,
	}, nil
}

func ScanDir(fn *DirNode, handler ScanningHandler) (*TreeNode, error) {
	curTreeNode := &TreeNode{
		info: fn,
	}

	entries, err := os.ReadDir(fn.absPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		// var curNode Node
		fileEntry, err := entry.Info()
		if err != nil {
			return nil, err
		}

		absEntryPath, err := filepath.Abs(fn.absPath + "/" + entry.Name())
		if err != nil {
			return nil, err
		}

		// TODO: Implement some common function (probably)
		if entry.IsDir() {
			dirNode := &DirNode{
				GeneralNode: GeneralNode{
					entry:   fileEntry,
					absPath: absEntryPath,
				},
			}

			err = handler.HandleDir(curTreeNode, dirNode)
			if err != nil {
				return nil, err
			}

			if childTreeNode, err := ScanDir(dirNode, handler); err != nil {
				return nil, err
			} else {
				curTreeNode.children = append(curTreeNode.children, childTreeNode)
			}
		} else {
			fileNode := &FileNode{
				GeneralNode: GeneralNode{
					entry:   fileEntry,
					absPath: absEntryPath,
				},
			}
			err = handler.HandleFile(curTreeNode, fileNode)
			if err != nil {
				return nil, err
			}

			if childTreeNode, err := ScanFile(fileNode, handler); err != nil {
				return nil, err
			} else {
				curTreeNode.children = append(curTreeNode.children, childTreeNode)
			}
		}
		// TODO: Convert this to BFS
	}

	return curTreeNode, nil
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
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return err
	}
	fn.absPath = obj.Parent
	fn.Tags = obj.Tags
	return nil
}

func (dn *DirNode) MarshalJSON() ([]byte, error) {
	obj := NodeJSON{
		Parent: dn.absPath,
		IsDir:  true,
		Tags:   dn.Tags,
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
	return nil
}
