package storage

import (
	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/cmderror"
	"github.com/heroku/self/MetaManager/internal/file"
	"github.com/heroku/self/MetaManager/internal/utils"
	"os"
	"path/filepath"
)

type FileStorageRW struct {
	dataFilePath string
}

func (f *FileStorageRW) Read() (*ds.TreeNode, error) {
	rootNode := ds.TreeNode{}
	rootNode.Serializer = file.FileNodeJSONSerializer{}

	// Read data in bytes from dataFilePath and construct TreeManager
	serializedNode, err := os.ReadFile(f.dataFilePath)
	if err != nil {
		return nil, err
	}

	// err = json.Unmarshal(serializedNode, &rootNode)
	err = rootNode.UnmarshalJSON(serializedNode)
	if err != nil {
		return nil, err
	}

	return &rootNode, nil
}

func (f *FileStorageRW) Write(root *ds.TreeNode) error {
	root.Serializer = file.FileNodeJSONSerializer{}

	// Save the tree structure in data.json
	// serializedNode, err := json.Marshal(root)
	// if err != nil {
	// 	return err
	// }

	serializedNode, err := root.MarshalJSON()
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
