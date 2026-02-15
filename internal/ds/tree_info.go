package ds

type InfoUnmarshaler interface {
	InfoUnmarshal(map[string]interface{}) (TreeNodeInformable, error)
}

type TreeNodeInformable interface {
	Name() string
}
