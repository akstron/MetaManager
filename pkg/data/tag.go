package data

import (
	"github/akstron/MetaManager/pkg/cmderror"
	"os"
)

/*
Tag related functionalities are implemented here
*/
type TagManager struct {
	dataFilePath string
	trMg         *TreeManager
}

func NewTagManager(dataFilePath string) (*TagManager, error) {
	tgMg := &TagManager{
		dataFilePath: dataFilePath,
	}

	// Read data in bytes from dataFilePath and construct TreeManager
	content, err := os.ReadFile(dataFilePath)
	if err != nil {
		return nil, err
	}

	/*
		WARNING: This does not initialize Root member of TreeManager
		There can be consequences
	*/
	tgMg.trMg = &TreeManager{}

	err = tgMg.trMg.Load(content)
	if err != nil {
		return nil, err
	}

	return tgMg, nil
}

func (tgMg *TagManager) AddTag(path string, tag string) error {
	treeNode, err := tgMg.trMg.FindNodeByAbsPath(path)
	if err != nil {
		return err
	}

	if fileNode, ok := treeNode.(*FileNode); ok {
		fileNode.Tags = append(fileNode.Tags, tag)
	} else if dirNode, ok := treeNode.(*DirNode); ok {
		dirNode.Tags = append(dirNode.Tags, tag)
	} else {
		return &cmderror.SomethingWentWrong{}
	}

	return nil
}

func (tgMg *TagManager) Save(dataFilePath string) error {
	return tgMg.trMg.Save(dataFilePath)
}
