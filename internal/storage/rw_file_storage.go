package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/heroku/self/MetaManager/internal/cmderror"
	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/file"
	"github.com/heroku/self/MetaManager/internal/utils"
)

type FileStorageRW struct {
	dataFilePath string
}

func buildTreeNodeFromJSON(jsonNode *ds.TreeNodeJSON, infoSerializer ds.InfoUnmarshaler) (*ds.TreeNode, error) {
	info, err := infoSerializer.InfoUnmarshal(jsonNode.Info)
	if err != nil {
		return nil, err
	}

	node := &ds.TreeNode{
		Info:     info,
		Children: []*ds.TreeNode{},
	}

	for _, child := range jsonNode.Children {
		childTreeNode, err := buildTreeNodeFromJSON(child, infoSerializer)
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, childTreeNode)
	}

	return node, nil
}

func (f *FileStorageRW) Read() (*ds.TreeNode, error) {
	serializedNode, err := os.ReadFile(f.dataFilePath)
	if err != nil {
		return nil, err
	}

	var rootNode ds.TreeNodeJSON

	err = json.Unmarshal(serializedNode, &rootNode)
	if err != nil {
		return nil, err
	}

	return buildTreeNodeFromJSON(&rootNode, &file.FileNodeJSONSerializer{})
}

func (f *FileStorageRW) Write(root *ds.TreeNode) error {
	serializedNode, err := json.Marshal(root)
	if err != nil {
		return err
	}

	err = os.WriteFile(f.dataFilePath, serializedNode, 0666)
	if err != nil {
		return err
	}

	return nil
}

func NewFileStorageRW(dataFilePath string) (*FileStorageRW, error) {
	return &FileStorageRW{
		dataFilePath: dataFilePath,
	}, nil
}

type FileStorageRWFactory struct {
	dirFilePath string
}

func (factory *FileStorageRWFactory) GetTreeRW() (TreeRW, error) {
	return NewFileStorageRW(factory.dirFilePath)
}

// GetRW returns a TreeRW for the given context's .mm directory. contextName must be non-empty and the .mm/<contextName> dir must exist.
func GetRW(contextName string) (TreeRW, error) {
	if contextName == "" {
		return nil, &cmderror.UninitializedRoot{}
	}
	found, root, err := utils.FindMMDirPath(contextName)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, &cmderror.UninitializedRoot{}
	}
	dataFilePath := filepath.Join(root, utils.DataFileName)
	return NewFileStorageRW(dataFilePath)
}
