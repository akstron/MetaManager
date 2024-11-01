package data

type TreeIterable interface {
	Next() (NodeInformable, error)
	HasNext() bool
}

type NodeIterable interface {
	GetFileChildren() []*FileNode
	GetDirChildren() []*DirNode
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

type TreeIterator struct {
	tgMg  *TreeManager
	index int
	nodes []NodeIterable
}

func (ti *TreeIterator) Next() (NodeInformable, error) {
	if ti.index >= len(ti.nodes) {
		return nil, nil
	}

	fileNodes := ti.nodes[ti.index].GetFileChildren()
	for _, filePtr := range fileNodes {
		if filePtr != nil {
			ti.nodes = append(ti.nodes, filePtr)
		}
	}

	dirNodes := ti.nodes[ti.index].GetDirChildren()
	for _, dirPtr := range dirNodes {
		if dirPtr != nil {
			ti.nodes = append(ti.nodes, dirPtr)
		}
	}

	ti.index++
	return ti.nodes[ti.index-1].GetInfoProvider(), nil
}

func (ti *TreeIterator) HasNext() bool {
	if ti.index >= len(ti.nodes) {
		return false
	}

	return true
}
