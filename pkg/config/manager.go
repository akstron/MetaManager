package config

import (
	"encoding/json"
	"fmt"
	"github/akstron/MetaManager/pkg/cmderror"
	"github/akstron/MetaManager/pkg/utils"
	"os"
	"path/filepath"
)

type Manager struct {
	DirPath string
	root    *Node
}

func (mg *Manager) ScanDirectory() error {
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
		return nil
	}

	rootNode := &DirNode{GeneralNode: GeneralNode{entry: root, absPath: dirPathAbs}}
	err = rootNode.Scan()

	if err != nil {
		return err
	}

	serializedNode, err := json.Marshal(rootNode)
	if err != nil {
		return err
	}

	var tempRootNode DirNode
	err = json.Unmarshal(serializedNode, &tempRootNode)
	fmt.Println(tempRootNode)

	if err != nil {
		fmt.Println(err)
	}

	return nil
}
