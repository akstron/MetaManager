package ds

/*
- Manages the nodes

DirPath
  - The parent root

root
  - Root of the scanned nodes. DirPath will be scanned
*/

type TreeManager struct {
	Root *TreeNode
}

func NewTreeManager(root *TreeNode) *TreeManager {
	return &TreeManager{
		Root: root,
	}
}
