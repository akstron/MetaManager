package data

import (
	"encoding/json"
	"github/akstron/MetaManager/pkg/cmderror"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

const (
	testDirPathRelative = "./testDir"
)

func createDirStructure() error {
	testDir, err := filepath.Abs(testDirPathRelative)
	if err != nil {
		return err
	}

	err = os.Mkdir(testDir, 0755)
	if err != nil {
		return err
	}

	return nil
}

func attachDir(root *DirNode, num int) *DirNode {
	dirNode := &DirNode{
		GeneralNode: GeneralNode{
			absPath: "dir" + strconv.Itoa(num),
		},
	}
	root.children = append(root.children, dirNode)
	return dirNode
}

func attachFile(root *DirNode, num int) *FileNode {
	fileNode := &FileNode{
		GeneralNode: GeneralNode{
			absPath: "file" + strconv.Itoa(num),
		},
	}
	root.children = append(root.children, fileNode)
	return fileNode
}

func dfs(root *DirNode, tree map[string][]string) error {
	for _, node := range root.children {
		if fileNode, ok := node.(*FileNode); ok {
			tree[root.absPath] = append(tree[root.absPath], fileNode.absPath)
		} else if dirNode, ok := node.(*DirNode); ok {
			tree[root.absPath] = append(tree[root.absPath], dirNode.absPath)
			dfs(dirNode, tree)
		} else {
			return &cmderror.SomethingWentWrong{}
		}
	}
	return nil
}

func isEqual(tree1, tree2 map[string][]string) bool {
	for key, value := range tree1 {
		tree2Val, ok := tree2[key]
		if !ok {
			return false
		}
		if len(value) != len(tree2Val) {
			return false
		}
		for i := 0; i < len(value); i++ {
			if value[i] != tree2Val[i] {
				return false
			}
		}
	}
	return true
}

func TestFileNodesSerialization(t *testing.T) {
	root := &DirNode{
		GeneralNode: GeneralNode{
			absPath: "dir0",
		},
	}
	dir1 := attachDir(root, 1)
	dir2 := attachDir(root, 2)
	dir3 := attachDir(dir1, 3)
	dir4 := attachDir(dir1, 4)
	dir5 := attachDir(dir3, 5)
	_ = attachFile(root, 1)
	_ = attachFile(root, 2)
	_ = attachFile(dir2, 3)
	_ = attachFile(dir5, 4)
	_ = attachFile(dir5, 5)
	_ = attachFile(dir5, 6)
	_ = attachFile(dir3, 7)
	attachFile(dir4, 8)

	var tree1 map[string][]string = map[string][]string{}
	err := dfs(root, tree1)
	if err != nil {
		t.Fatal(err)
	}

	serializedTree, err := json.Marshal(root)
	if err != nil {
		t.Fatal(err)
	}

	var root2 *DirNode
	err = json.Unmarshal(serializedTree, &root2)
	if err != nil {
		t.Fatal(err)
	}

	tree2 := make(map[string][]string)
	err = dfs(root2, tree2)
	if err != nil {
		t.Fatal(err)
	}

	if !isEqual(tree1, tree2) || !isEqual(tree2, tree1) {
		t.FailNow()
	}
}

// func TestFileStructureCreation(t *testing.T) {
// 	var dirPath string
// 	err := createDirStructure()
// 	if err != nil {
// 		goto finally
// 	}

// 	dirPath, err = filepath.Abs(testDirPathRelative)
// 	if err != nil {
// 		goto finally
// 	}

// 	InitRoot(dirPath)

// finally:
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }
