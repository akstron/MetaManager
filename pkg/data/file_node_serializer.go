package data

import "encoding/json"

type NodeJSON struct {
	Parent   string
	Children [][]byte
	IsDir    bool
	Tags     []string
}

func (fn *FileNode) MarshalJSON() ([]byte, error) {
	obj := NodeJSON{
		Parent: fn.absPath,
		IsDir:  false,
		Tags:   fn.Tags,
	}
	return json.Marshal(obj)
}

func (fn *FileNode) UnmarshalJSON(data []byte) error {
	var obj NodeJSON
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return err
	}
	fn.absPath = obj.Parent
	fn.Tags = obj.Tags
	return nil
}

func (dn *DirNode) MarshalJSON() ([]byte, error) {
	obj := NodeJSON{
		Parent: dn.absPath,
		IsDir:  true,
		Tags:   dn.Tags,
	}

	return json.Marshal(&obj)
}

func (dn *DirNode) UnmarshalJSON(data []byte) error {
	var obj NodeJSON
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return err
	}
	dn.absPath = obj.Parent
	dn.Tags = obj.Tags
	return nil
}
