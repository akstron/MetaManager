package data

import (
	"fmt"
	"github/akstron/MetaManager/pkg/cmderror"
)

/*
Tag related functionalities are implemented here
*/
type TagManager struct {
	trMg *TreeManager
}

func NewTagManager() *TagManager {
	return &TagManager{}
}

/*
TODO: Create a TagReader interface instead
This way we can decouple tree reading writing from tag
probably
*/
func (tgMg *TagManager) Load(r TreeReader) error {
	var err error

	tgMg.trMg = &TreeManager{}
	tgMg.trMg.Root, err = r.Read()
	if err != nil {
		return err
	}
	return nil
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

func (tgMg *TagManager) GetTag() ([]string, error) {
	if tgMg.trMg == nil {
		return nil, fmt.Errorf("Invalid operation. Tree not loaded")
	}
	return nil, nil
}

func (tgMg *TagManager) Save(rw TreeRW) error {
	return rw.Write(tgMg.trMg.Root)
}
