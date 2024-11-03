package ds

import "encoding/json"

func (tn *TreeNode) UnmarshalJSON(data []byte) error {
	var obj TreeNodeJSON
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return err
	}

	tn.Info, err = tn.Serializer.InfoUnmarshal(obj.Info, obj.SerializationInfo)
	if err != nil {
		return err
	}

	for _, child := range obj.Children {
		var childNode TreeNode
		childNode.Serializer = tn.Serializer
		// err := json.Unmarshal(child, &childNode)
		err := childNode.UnmarshalJSON(child)
		if err != nil {
			return err
		}
		tn.Children = append(tn.Children, &childNode)
	}

	return nil
}

func (tn *TreeNode) MarshalJSON() ([]byte, error) {
	var infoSerialized []byte

	infoSerialized, serializationInfo, err := tn.Serializer.InfoMarshal(tn.Info)
	if err != nil {
		return nil, err
	}

	var childrenSerialized [][]byte
	for _, node := range tn.Children {
		node.Serializer = tn.Serializer
		// currentChildSerialized, err := json.Marshal(node)
		currentChildSerialized, err := node.MarshalJSON()
		if err != nil {
			return nil, err
		}
		childrenSerialized = append(childrenSerialized, (currentChildSerialized))
	}

	obj := TreeNodeJSON{
		Info:              infoSerialized,
		Children:          childrenSerialized,
		SerializationInfo: serializationInfo,
	}

	return json.Marshal(&obj)
}
