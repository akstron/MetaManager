package ds

type InfoSerializer interface {
	// bytes, SerializationInfo, error
	InfoMarshal(TreeNodeInformable) ([]byte, string, error)
	InfoUnmarshal([]byte, string) (TreeNodeInformable, error)
}

type TreeNodeInformable interface {
	Name() string
}

/*
TreeManager is the owner of the dir structure.
All tree related operations should happen with the help
of TreeManager
*/
type TreeNode struct {
	/*Store any info in a node*/
	Info       TreeNodeInformable
	Children   []*TreeNode
	Serializer InfoSerializer
}

type TreeNodeJSON struct {
	Info              []byte
	Children          [][]byte
	SerializationInfo string
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
