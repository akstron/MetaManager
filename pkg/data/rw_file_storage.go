package data

import (
	"encoding/json"
	"os"
)

type FileStorageRW struct {
	dataFilePath string
}

func (f *FileStorageRW) Read() (*TreeNode, error) {
	rootNode := &TreeNode{}

	// Read data in bytes from dataFilePath and construct TreeManager
	serializedNode, err := os.ReadFile(f.dataFilePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(serializedNode, rootNode)
	if err != nil {
		return nil, err
	}

	return rootNode, nil
}

func (f *FileStorageRW) Write(root *TreeNode) error {
	// Save the tree structure in data.json
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
