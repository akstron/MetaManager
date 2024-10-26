package data

/*
	In future, we might want to change the backend storage
	system, that's why we have an interface
*/

// Fine-grained interface so that we can optimally read trees according
// to the storage system
type TreeReader interface {
	Read(*DirNode) error
}

// Fine-grained interface so that we can optimally write trees according
// to the storage system
type TreeWriter interface {
	Write(*DirNode) error
}

type TreeRW interface {
	TreeReader
	TreeWriter
}

func WriteTree(tw TreeWriter, node *DirNode) error {
	return tw.Write(node)
}

func ReadTree(tr TreeReader, node *DirNode) error {
	return tr.Read(node)
}
