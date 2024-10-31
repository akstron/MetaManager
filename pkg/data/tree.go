package data

import (
	"fmt"
	"github/akstron/MetaManager/pkg/cmderror"
	"github/akstron/MetaManager/pkg/config"
	"github/akstron/MetaManager/pkg/utils"
	"os"
	"path/filepath"
)

/*
- Manages the nodes

DirPath
  - The parent root

root
  - Root of the scanned nodes. DirPath will be scanned
*/

/*
TreeManager is the owner of the dir structure.
All tree related operations should happen with the help
of TreeManager
*/
type TreeManager struct {
	DirPath string
	Root    *DirNode
}

func NewTreeManager() TreeManager {
	return TreeManager{}
}

type NodeAbsPathIgnorer struct {
	igMg *config.IgnoreManager
}

func (ig *NodeAbsPathIgnorer) ShouldIgnore(ignorePath string) (bool, error) {
	/*
		igMg can have a GetData which returns constant data for iteration
		but lets see if this should be done
	*/
	for _, value := range ig.igMg.Data.Paths {
		if value == ignorePath {
			return true, nil
		}
	}
	return false, nil
}

func (mg *TreeManager) FindNodeByAbsPath(path string) (any, error) {
	// TODO -> Optimization probably
	queue := []*DirNode{}
	queue = append(queue, mg.Root)
	for i := 0; i < len(queue); i++ {
		if queue[i].absPath == path {
			return queue[i], nil
		}
		for _, child := range queue[i].FileChildren {
			if child.absPath == path {
				return child, nil
			}
		}
		queue = append(queue, queue[i].DirChildren...)
	}
	return nil, nil
}

/*
TODO: Remove this from TreeManager as this shouldn't be
aware of directory structure. Need to think about decoupling it.
*/
func (mg *TreeManager) ScanDirectory() error {
	present, err := utils.IsFilePresent(mg.DirPath)
	if err != nil {
		return err
	}

	if !present {
		return &cmderror.InvalidPath{}
	}

	dirPathAbs, err := filepath.Abs(mg.DirPath)
	if err != nil {
		return err
	}

	fmt.Println(dirPathAbs)

	root, err := os.Stat(dirPathAbs)
	if err != nil {
		return err
	}

	mg.Root = &DirNode{GeneralNode: GeneralNode{entry: root, absPath: dirPathAbs}}

	igMg, err := config.NewIgnoreManager()
	if err != nil {
		return err
	}

	igMg.Load()
	ignorer := &NodeAbsPathIgnorer{
		igMg: igMg,
	}

	/*
		TODO: Probably better to pass pointer
		We will use go routines to accelerate dfs
	*/
	handler := ScanHandler{
		ig: ignorer,
	}
	err = mg.Root.Scan(handler)

	if err != nil {
		return err
	}

	return nil
}
