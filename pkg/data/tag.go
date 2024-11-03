package data

import (
	"fmt"
	"github/akstron/MetaManager/ds"
	"github/akstron/MetaManager/pkg/file"
	"github/akstron/MetaManager/storage"
)

/*
Tag related functionalities are implemented here
*/
type TagManager struct {
	trMg *DirTreeManager
}

func NewTagManager() *TagManager {
	return &TagManager{}
}

/*
TODO: Create a TagReader interface instead
This way we can decouple tree reading writing from tag
probably
*/
func (tgMg *TagManager) Load(r storage.TreeReader) error {
	var err error

	root, err := r.Read()
	if err != nil {
		return err
	}

	tgMg.trMg = &DirTreeManager{
		TreeManager: ds.NewTreeManager(root),
	}
	return nil
}

func (tgMg *TagManager) AddTag(path string, tag string) error {
	nodeInfo, err := tgMg.trMg.FindNodeByAbsPath(path)
	if err != nil {
		return err
	}

	if nodeInfo == nil {
		return fmt.Errorf("path: %s not tracked", path)
	}

	nodeInfo.AddTag(tag)

	return nil
}

func (tgMg *TagManager) GetTaggedNodes(tag string) ([]string, error) {
	if tgMg.trMg == nil {
		return nil, fmt.Errorf("invalid operation, tree not loaded")
	}

	it := ds.NewTreeIterator(tgMg.trMg.TreeManager)
	return tgMg.iterateAndExtractPathsWithTag(&it, tag)
}

func (tgMg *TagManager) GetNodeTags(path string) ([]string, error) {
	if tgMg.trMg == nil {
		return nil, fmt.Errorf("invalid operation, tree not loaded")
	}

	nodeInfo, err := tgMg.trMg.FindNodeByAbsPath(path)
	if err != nil {
		return nil, err
	}

	return nodeInfo.GetTags(), nil
}

func IsPresent(val string, container []string) bool {
	for _, eachVal := range container {
		if eachVal == val {
			return true
		}
	}
	return false
}

func (*TagManager) iterateAndExtractPathsWithTag(it ds.TreeIterable, tag string) ([]string, error) {
	result := []string{}
	for it.HasNext() {
		curNode, err := it.Next()
		got := curNode.Info
		if err != nil {
			return nil, err
		}

		nodeTags := got.(file.NodeInformable).GetTags()
		if IsPresent(tag, nodeTags) {
			result = append(result, got.(file.NodeInformable).GetAbsPath())
		}
	}
	return result, nil
}

func (tgMg *TagManager) Save(rw storage.TreeRW) error {
	return rw.Write(tgMg.trMg.Root)
}
