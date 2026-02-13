package file

import (
	"encoding/json"
	"fmt"

	"github.com/heroku/self/MetaManager/internal/ds"
)

// NodeJSON is the persisted shape for all nodes (local and gdrive, file and dir).
type NodeJSON struct {
	Parent  string   `json:"Parent"`
	Tags    []string `json:"Tags"`
	Id      string   `json:"Id"`
	DriveId string   `json:"DriveId,omitempty"`
}

type FileNodeJSONSerializer struct{}

func (FileNodeJSONSerializer) InfoMarshal(info ds.TreeNodeInformable) ([]byte, string, error) {
	if _, ok := info.(*FileNode); !ok {
		return nil, "", fmt.Errorf("expected *FileNode, got %T", info)
	}
	serializedInfo, err := json.Marshal(info)
	return serializedInfo, "FILE", err
}

func (FileNodeJSONSerializer) InfoUnmarshal(data []byte, serializationInfo string) (ds.TreeNodeInformable, error) {
	var fn FileNode
	if err := json.Unmarshal(data, &fn); err != nil {
		return nil, err
	}
	return &fn, nil
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
