package file

import (
	"encoding/json"
	"io/fs"
)

/*Common node operations which should be provided by all nodes*/
type NodeInformable interface {
	GetAbsPath() string
	GetTags() []string
	AddTag(string)
	DeleteTag(string)
	SetId(string)
	GetId() string
}
type GeneralNode struct {
	AbsPath string      `mapstructure:"Parent"`
	Entry   fs.FileInfo `mapstructure:"Entry"`
	Tags    []string    `mapstructure:"Tags"`
	// User friendly id, which uniquely finds a node
	// exception: empty string
	Id string `mapstructure:"Id"`
}

func NewGeneralNode(absPath string, entry fs.FileInfo) GeneralNode {
	return GeneralNode{
		AbsPath: absPath,
		Entry:   entry,
	}
}

func (gn *GeneralNode) GetAbsPath() string {
	return gn.AbsPath
}

func (gn *GeneralNode) GetTags() []string {
	return gn.Tags
}

func (gn *GeneralNode) AddTag(tag string) {
	for _, lsTag := range gn.Tags {
		if lsTag == tag {
			return
		}
	}

	gn.Tags = append(gn.Tags, tag)
}

func (gn *GeneralNode) DeleteTag(tag string) {
	newTagList := []string{}
	for _, tagLs := range gn.Tags {
		if tagLs != tag {
			newTagList = append(newTagList, tagLs)
		}
	}
	gn.Tags = newTagList
}

func (gn *GeneralNode) SetId(id string) {
	gn.Id = id
}

func (gn *GeneralNode) GetId() string {
	return gn.Id
}

type SerializableNode interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}

// FileNode represents a file or directory (local or gdrive). DriveId is set for Google Drive nodes.
type FileNode struct {
	GeneralNode `mapstructure:",squash"`
	DriveId     string `mapstructure:"DriveId"` // non-empty for Google Drive nodes
}

func (fn *FileNode) GetInfoProvider() NodeInformable {
	return fn
}

func (fn *FileNode) Name() string {
	return "FILE"
}

// NodeJSON is the JSON representation of FileNode, using "Parent" instead of "AbsPath"
type NodeJSON struct {
	Parent  string   `json:"Parent"`
	Tags    []string `json:"Tags,omitempty"`
	Id      string   `json:"Id,omitempty"`
	DriveId string   `json:"DriveId,omitempty"`
}

func (fn *FileNode) MarshalJSON() ([]byte, error) {
	obj := NodeJSON{
		Parent:  fn.AbsPath,
		Tags:    fn.Tags,
		Id:      fn.Id,
		DriveId: fn.DriveId,
	}
	return json.Marshal(obj)
}

func (fn *FileNode) UnmarshalJSON(data []byte) error {
	var obj NodeJSON
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	fn.AbsPath = obj.Parent
	fn.Tags = obj.Tags
	fn.Id = obj.Id
	fn.DriveId = obj.DriveId
	return nil
}
