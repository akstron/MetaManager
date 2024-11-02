package data

/*
TreeManager is the owner of the dir structure.
All tree related operations should happen with the help
of TreeManager
*/
type TreeNode struct {
	/*Store any info in a node*/
	info     any
	children []*TreeNode
}

type TreeNodeJSON struct {
	Info              []byte
	Children          [][]byte
	SerializationInfo string
}

func NewTreeNode(info any) TreeNode {
	return TreeNode{
		info: info,
	}
}

func (tn *TreeNode) GetChildren() []*TreeNode {
	return tn.children
}
