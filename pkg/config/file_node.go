package config

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
)

type GeneralNode struct {
	absPath string
	entry   fs.FileInfo
}

type ScannableNode interface {
	Scan() error
}

type SerializableNode interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}

type Node interface {
	/*
		Composition design pattern
		Both FileNode and DirNode implements this interface
	*/
	ScannableNode
	// All nodes should be able to tell how they can be serialized/derialized
	SerializableNode
}

type FileNode struct {
	GeneralNode
}

func (fn *FileNode) Scan() error {
	return nil
}

type DirNode struct {
	GeneralNode
	children []Node
}

func (fn *DirNode) Scan() error {
	entries, err := os.ReadDir(fn.absPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		var curNode Node
		fileEntry, err := entry.Info()
		if err != nil {
			return err
		}

		absEntryPath, err := filepath.Abs(fn.absPath + "/" + entry.Name())
		if err != nil {
			return err
		}

		// TODO: Implement factory pattern
		if entry.IsDir() {
			curNode = &DirNode{
				GeneralNode: GeneralNode{
					entry:   fileEntry,
					absPath: absEntryPath,
				},
			}
		} else {
			curNode = &FileNode{
				GeneralNode: GeneralNode{
					entry:   fileEntry,
					absPath: absEntryPath,
				},
			}
		}
		fn.children = append(fn.children, curNode)
		// TODO: Convert this to BFS
		if err := curNode.Scan(); err != nil {
			return err
		}
	}

	return nil
}

type NodeJSON struct {
	Parent   string
	Children [][]byte
	IsDir    bool
}

func (fn *FileNode) MarshalJSON() ([]byte, error) {
	obj := NodeJSON{
		Parent: fn.absPath,
		IsDir:  false,
	}
	return json.Marshal(obj)
}

func (fn *FileNode) UnmarshalJSON(data []byte) error {
	var obj NodeJSON
	return json.Unmarshal(data, &obj)
}

func (dn *DirNode) MarshalJSON() ([]byte, error) {
	var childrenSerialized [][]byte
	for _, node := range dn.children {
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
	for _, child := range obj.Children {
		var childObj NodeJSON
		err := json.Unmarshal(child, &childObj)
		if err != nil {
			return err
		}

		childAbsPath := childObj.Parent
		var curNode Node
		if obj.IsDir {
			curNode = &DirNode{
				GeneralNode: GeneralNode{
					absPath: childAbsPath,
				},
			}
		} else {
			curNode = &FileNode{
				GeneralNode: GeneralNode{
					absPath: childAbsPath,
				},
			}
		}
		dn.children = append(dn.children, curNode)
		err = json.Unmarshal(child, &curNode)
		if err != nil {
			return err
		}
	}

	return nil
}
