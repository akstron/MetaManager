package file

import (
	"encoding/json"
	"fmt"
	"github.com/heroku/self/MetaManager/internal/ds"
)

type NodeJSON struct {
	Parent string
	Tags   []string
	Id     string
}

type FileNodeJSONSerializer struct {
}

func (FileNodeJSONSerializer) InfoMarshal(info ds.TreeNodeInformable) ([]byte, string, error) {
	fileType := "DIR"
	switch info.(type) {
	case *FileNode:
		fileType = "FILE"
	case *DirNode:
		fileType = "DIR"
	case *DriveFileNode:
		fileType = "GDRIVE_FILE"
	case *DriveDirNode:
		fileType = "GDRIVE_DIR"
	}
	serializedInfo, err := json.Marshal(info)
	return serializedInfo, fileType, err
}

func (FileNodeJSONSerializer) InfoUnmarshal(data []byte, serializationInfo string) (ds.TreeNodeInformable, error) {
	var result ds.TreeNodeInformable
	switch serializationInfo {
	case "FILE":
		var fileNode FileNode
		err := json.Unmarshal(data, &fileNode)
		if err != nil {
			return nil, err
		}
		result = &fileNode
	case "DIR":
		var dirNode DirNode
		err := json.Unmarshal(data, &dirNode)
		if err != nil {
			return nil, err
		}
		result = &dirNode
	case "GDRIVE_FILE":
		var n DriveFileNode
		err := json.Unmarshal(data, &n)
		if err != nil {
			return nil, err
		}
		result = &n
	case "GDRIVE_DIR":
		var n DriveDirNode
		err := json.Unmarshal(data, &n)
		if err != nil {
			return nil, err
		}
		result = &n
	default:
		return nil, fmt.Errorf("unknown serializationInfo: %s found", serializationInfo)
	}
	return result, nil
}

func (fn *FileNode) MarshalJSON() ([]byte, error) {
	obj := NodeJSON{
		Parent: fn.AbsPath,
		Tags:   fn.Tags,
		Id:     fn.Id,
	}
	return json.Marshal(obj)
}

func (fn *FileNode) UnmarshalJSON(data []byte) error {
	var obj NodeJSON
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return err
	}
	fn.AbsPath = obj.Parent
	fn.Tags = obj.Tags
	fn.Id = obj.Id
	return nil
}

func (dn *DirNode) MarshalJSON() ([]byte, error) {
	obj := NodeJSON{
		Parent: dn.AbsPath,
		Tags:   dn.Tags,
		Id:     dn.Id,
	}

	return json.Marshal(&obj)
}

func (dn *DirNode) UnmarshalJSON(data []byte) error {
	var obj NodeJSON
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return err
	}
	dn.AbsPath = obj.Parent
	dn.Tags = obj.Tags
	dn.Id = obj.Id
	return nil
}
