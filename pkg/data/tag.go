package data

import (
	"github/akstron/MetaManager/pkg/cmderror"
)

/*
Tag related functionalities are implemented here
*/
type TagManager struct {
	trMg *TreeManager
}

func NewTagManager(rw TreeReader) (*TagManager, error) {
	var err error

	tgMg := &TagManager{}
	tgMg.trMg = &TreeManager{}
	tgMg.trMg.Root, err = rw.Read()
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

func (tgMg *TagManager) Save(rw TreeRW) error {
	return rw.Write(tgMg.trMg.Root)
}
