package data

import (
	"encoding/json"
)

func (tn *TreeNode) UnmarshalJSON(data []byte) error {
	var obj TreeNodeJSON
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return err
	}

	if obj.SerializationInfo == "FILE" {
		var fileNode FileNode
		err := json.Unmarshal(obj.Info, &fileNode)
		if err != nil {
			return nil
		}
		tn.info = &fileNode
	} else {
		var dirNode DirNode
		err = json.Unmarshal(obj.Info, &dirNode)
		if err != nil {
			return nil
		}
		tn.info = &dirNode
	}

	for _, child := range obj.Children {
		var childNode TreeNode
		err := json.Unmarshal(child, &childNode)
		if err != nil {
			return err
		}
		tn.children = append(tn.children, &childNode)
	}

	return nil
}

func (tn *TreeNode) MarshalJSON() ([]byte, error) {
	var infoSerialized []byte
	infoSerialized, err := json.Marshal(tn.info)
	if err != nil {
		return nil, err
	}

	var childrenSerialized [][]byte
	for _, node := range tn.children {
		currentChildSerialized, err := json.Marshal(node)
		if err != nil {
			return nil, err
		}
		childrenSerialized = append(childrenSerialized, (currentChildSerialized))
	}

	_, ok := tn.info.(*FileNode)
	fileType := "DIR"
	if ok {
		fileType = "FILE"
	}

	obj := TreeNodeJSON{
		Info:              infoSerialized,
		Children:          childrenSerialized,
		SerializationInfo: fileType,
	}

	return json.Marshal(&obj)
}
