package ds

/*
TreeManager is the owner of the dir structure.
All tree related operations should happen with the help
of TreeManager
*/
type TreeNode struct {
	/*Store any info in a node*/
	Info     TreeNodeInformable `json:"info"`
	Children []*TreeNode        `json:"children"`
}

type TreeNodeJSON struct {
	Info     map[string]interface{} `json:"info"`
	Children []*TreeNodeJSON        `json:"children"`
}

func NewTreeNode(info TreeNodeInformable) *TreeNode {
	return &TreeNode{
		Info: info,
	}
}

func (tn *TreeNode) GetChildren() []*TreeNode {
	return tn.Children
}

func (tn *TreeNode) AddChild(node *TreeNode) {
	tn.Children = append(tn.Children, node)
}
