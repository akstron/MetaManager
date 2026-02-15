package file

import (
	mapstructure "github.com/go-viper/mapstructure/v2"
	"github.com/heroku/self/MetaManager/internal/ds"
)

type FileNodeJSONSerializer struct{}

func (FileNodeJSONSerializer) InfoUnmarshal(info map[string]interface{}) (ds.TreeNodeInformable, error) {
	var fn FileNode
	err := mapstructure.Decode(info, &fn)
	if err != nil {
		return nil, err
	}
	return &fn, nil
}

// Fail build if FileNodeJSONSerializer does not implement InfoUnmarshaler
var _ ds.InfoUnmarshaler = (*FileNodeJSONSerializer)(nil)
