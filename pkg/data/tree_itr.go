package data

type TreeIterable interface {
	Next() (any, error)
	HasNext() (bool, error)
}

type NodeIterable interface {
	GetFileChildren() []*FileNode
	GetDirChildren() []*DirNode
}

func NewTreeIterator(tgMg TreeManager) TreeIterator {
	tI := TreeIterator{
		tgMg: tgMg,
	}

	tI.index = 0
	tI.nodes = append(tI.nodes, tI.tgMg.Root)

	return tI
}

type TreeIterator struct {
	tgMg  TreeManager
	index int
	nodes []NodeIterable
}

func (ti *TreeIterator) Next() (any, error) {
	if ti.index >= len(ti.nodes) {
		return false, nil
	}

	fileNodes := ti.nodes[ti.index].GetFileChildren()
	if fileNodes != nil {
		for _, filePtr := range fileNodes {
			ti.nodes = append(ti.nodes, filePtr)
		}
	}

	dirNodes := ti.nodes[ti.index].GetDirChildren()
	if dirNodes != nil {
		for _, dirPtr := range dirNodes {
			ti.nodes = append(ti.nodes, dirPtr)
		}
	}

	ti.index++
	return ti.nodes[ti.index-1], nil
}

func (ti *TreeIterator) HasNext() bool {
	if ti.index >= len(ti.nodes) {
		return false
	}

	return true
}
