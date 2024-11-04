package utils

import (
	"github/akstron/MetaManager/ds"
	"github/akstron/MetaManager/pkg/file"
	"strings"

	"github.com/jedib0t/go-pretty/v6/list"
)

func ConstructTreeWriter(curNode *ds.TreeNode, cutPrefix string, wr list.Writer) error {
	info := curNode.Info.(file.NodeInformable)
	insPath, _ := strings.CutPrefix(info.GetAbsPath(), cutPrefix+"/")
	wr.AppendItem(insPath)

	wr.Indent()

	for _, child := range curNode.Children {
		err := ConstructTreeWriter(child, info.GetAbsPath(), wr)
		if err != nil {
			return err
		}
	}

	wr.UnIndent()
	return nil
}
