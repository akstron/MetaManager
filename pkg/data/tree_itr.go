package data

type TreeIterable interface {
	Next() (NodeInformable, error)
	HasNext() bool
}

type NodeIterable interface {
	GetChildren() []NodeIterable
	GetInfoProvider() NodeInformable
}

func NewTreeIterator(tgMg *TreeManager) TreeIterator {
	tI := TreeIterator{
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
	nodes []NodeIterable
}

func (ti *TreeIterator) Next() (NodeInformable, error) {
	if ti.index >= len(ti.nodes) {
		return nil, nil
	}

	childNodes := ti.nodes[ti.index].GetChildren()
	for _, childNode := range childNodes {
		if childNode != nil {
			ti.nodes = append(ti.nodes, childNode)
		}
	}

	ti.index++
	return ti.nodes[ti.index-1].GetInfoProvider(), nil
}

func (ti *TreeIterator) HasNext() bool {
	return ti.index < len(ti.nodes)
}
