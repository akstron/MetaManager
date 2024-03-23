package config

import (
	"encoding/json"
	"fmt"
	"github/akstron/MetaManager/pkg/cmderror"
	"github/akstron/MetaManager/pkg/utils"
	"io/fs"
	"os"
	"path/filepath"
)

type GeneralNode struct {
	absPath string
	entry   fs.FileInfo
}

type Node interface {
	Scan() error
}
type FileNode struct {
	GeneralNode
}

func (fn *FileNode) Scan() error {
	return nil
}

type DirNode struct {
	GeneralNode
	children []Node
}

func (fn *DirNode) Scan() error {
	// absPath, err := filepath.Abs("./" + fn.entry.Name())
	// fmt.Println(absPath)
	// if err != nil {
	// 	return err
	// }

	entries, err := os.ReadDir(fn.absPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		var curNode Node
		fileEntry, err := entry.Info()
		if err != nil {
			return err
		}

		fmt.Println(fileEntry.Name(), "->")
		absEntryPath, err := filepath.Abs(fn.absPath + "/" + entry.Name())
		if err != nil {
			return err
		}

		if entry.IsDir() {
			curNode = &DirNode{
				GeneralNode: GeneralNode{
					entry:   fileEntry,
					absPath: absEntryPath,
				},
			}
		} else {
			curNode = &FileNode{
				GeneralNode: GeneralNode{
					entry:   fileEntry,
					absPath: absEntryPath,
				},
			}
		}
		fn.children = append(fn.children, curNode)
		// TODO: Convert this to BFS
		if err := curNode.Scan(); err != nil {
			return err
		}

		fmt.Println()
	}

	return nil
}

func ScanDirectory(dirPath string) error {
	present, err := utils.IsFilePresent(dirPath)
	if err != nil {
		return err
	}

	if !present {
		return &cmderror.InvalidPath{}
	}

	dirPathAbs, err := filepath.Abs(dirPath)
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

	res, _ := json.Marshal(rootNode)
	fmt.Println(res)

	return nil
}
