package ds

// import (
// 	"testing"
// )

// func isFileSame(x any, y *FileNode) (bool, error) {
// 	fileNode, ok := x.(*FileNode)
// 	if !ok {
// 		return false, nil
// 	}
// 	return (fileNode == y), nil
// }

// func isDirSame(x any, y *DirNode) (bool, error) {
// 	dirNode, ok := x.(*DirNode)
// 	if !ok {
// 		return false, nil
// 	}
// 	return (dirNode == y), nil
// }

// func CheckDir(t *testing.T, trMg *TreeManager, absPath string, nodeToCompare *DirNode) {
// 	node, err := trMg.FindNodeByAbsPath(absPath)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	isEqual, err := isDirSame(node, nodeToCompare)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if !((isEqual && absPath == nodeToCompare.absPath) || (!isEqual && absPath != nodeToCompare.absPath)) {
// 		t.FailNow()
// 	}
// }

// func CheckFile(t *testing.T, trMg *TreeManager, absPath string, nodeToCompare *FileNode) {
// 	node, err := trMg.FindNodeByAbsPath(absPath)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	isEqual, err := isFileSame(node, nodeToCompare)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if !((isEqual && absPath == nodeToCompare.absPath) || (!isEqual && absPath != nodeToCompare.absPath)) {
// 		t.FailNow()
// 	}
// }

// func TestFindNodeByAbsPath(t *testing.T) {
// 	root := &DirNode{
// 		GeneralNode: GeneralNode{
// 			absPath: "dir0",
// 		},
// 	}
// 	dir1 := attachDir(root, 1)
// 	dir2 := attachDir(root, 2)
// 	dir3 := attachDir(dir1, 3)
// 	dir4 := attachDir(dir1, 4)
// 	dir5 := attachDir(dir3, 5)
// 	file1 := attachFile(root, 1)
// 	file2 := attachFile(root, 2)
// 	_ = attachFile(dir2, 3)
// 	_ = attachFile(dir5, 4)
// 	_ = attachFile(dir5, 5)
// 	_ = attachFile(dir5, 6)
// 	_ = attachFile(dir3, 7)
// 	attachFile(dir4, 8)

// 	trMg := &TreeManager{
// 		Root: root,
// 	}

// 	CheckDir(t, trMg, "dir0", root)
// 	CheckDir(t, trMg, "dirx", root)
// 	CheckDir(t, trMg, "dir1", dir3)
// 	CheckDir(t, trMg, "dir1", dir1)
// 	CheckDir(t, trMg, "dir1", dir2)
// 	CheckDir(t, trMg, "dir2", dir2)

// 	CheckFile(t, trMg, "file1", file1)
// 	CheckFile(t, trMg, "file2", file2)
// 	CheckFile(t, trMg, "file3", file1)
// }
