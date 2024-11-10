package ds

type TreeIterable interface {
	Next() (*TreeNode, error)
	HasNext() bool
}

type NodeIterable interface {
	GetChildren() []NodeIterable
}

func NewTreeIterator(tgMg *TreeManager) *TreeIterator {
	tI := &TreeIterator{
		tgMg: tgMg,
	}

	tI.index = 0
	tI.nodes = append(tI.nodes, tI.tgMg.Root)

	return tI
}

/*
The iterator iterates over all the nodes in the tree managed by TreeManager
It does not matters if the TreeManager manages a subtree of a bigger tree
*/
type TreeIterator struct {
	tgMg  *TreeManager
	index int
	nodes []*TreeNode
}

/*
TODO: Change this to return any.
*/
func (ti *TreeIterator) Next() (*TreeNode, error) {
	if ti.index >= len(ti.nodes) {
		return nil, nil
	}

	childNodes := ti.nodes[ti.index].GetChildren()
	for _, childNode := range childNodes {
		if childNode != nil && childNode.Info != nil {
			ti.nodes = append(ti.nodes, childNode)
		}
	}

	ti.index++
	data := ti.nodes[ti.index-1]
	return data, nil
}

func (ti *TreeIterator) HasNext() bool {
	return ti.index < len(ti.nodes)
}
