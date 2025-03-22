package storage

import (
	"github/akstron/MetaManager/ds"
	"github/akstron/MetaManager/pkg/cmderror"
	"github/akstron/MetaManager/pkg/file"
	"github/akstron/MetaManager/pkg/utils"
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

/*
TODO: This will be changed based on certain flags -> Currently not implemented
Example: user can set if it wants to use filesystem or some database implementation
*/
func GetRW() (TreeRW, error) {
	found, root, err := utils.FindMMDirPath()
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, &cmderror.UninitializedRoot{}
	}

	/*
		Check if the data.json is already written.
		Don't override, if already written
	*/
	dataFilePath := filepath.Join(root, utils.DATA_FILE_NAME)

	return NewFileStorageRW(dataFilePath)
}
