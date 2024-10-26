package data

import (
	"encoding/json"
	"os"
)

type FileStorageRW struct {
	DataFilePath string
}

func (f *FileStorageRW) Read(root *DirNode) error {
	return nil
}

func (f *FileStorageRW) Write(root *DirNode) error {
	// Save the tree structure in data.json
	serializedNode, err := json.Marshal(root)
	if err != nil {
		return err
	}

	err = os.WriteFile(f.DataFilePath, serializedNode, 0666)
	if err != nil {
		return err
	}

	return nil
}

func NewFileStorageRW(dataFilePath string) (*FileStorageRW, error) {
	return &FileStorageRW{
		DataFilePath: dataFilePath,
	}, nil
}
