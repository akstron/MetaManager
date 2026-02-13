package file

import (
	"encoding/json"
	"testing"

	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/stretchr/testify/require"
)

func TestDataSerialization(t *testing.T) {
	root := &ds.TreeNode{
		Info: &FileNode{GeneralNode: GeneralNode{AbsPath: "/"}},
		Children: []*ds.TreeNode{
			{
				Info: &FileNode{GeneralNode: GeneralNode{AbsPath: "/f"}},
			},
			{
				Info: &FileNode{GeneralNode: GeneralNode{AbsPath: "/d"}},
				Children: []*ds.TreeNode{
					{
						Info: &FileNode{GeneralNode: GeneralNode{AbsPath: "/d/a"}},
						Children: []*ds.TreeNode{
							{
								Info: &FileNode{GeneralNode: GeneralNode{AbsPath: "/d/a/x"}},
							},
						},
					},
					{
						Info: &FileNode{GeneralNode: GeneralNode{AbsPath: "/d/b"}},
					},
					{
						Info: &FileNode{GeneralNode: GeneralNode{AbsPath: "/d/c"}},
					},
					{
						Info: &FileNode{GeneralNode: GeneralNode{AbsPath: "/d/g"}},
					},
				},
			},
			{
				Info: &FileNode{GeneralNode: GeneralNode{AbsPath: "/h"}},
			},
		},
	}

	root.Serializer = FileNodeJSONSerializer{}
	serializedNode, err := json.Marshal(root)
	require.NoError(t, err)

	var extractedRoot ds.TreeNode
	extractedRoot.Serializer = FileNodeJSONSerializer{}
	err = json.Unmarshal(serializedNode, &extractedRoot)
	require.NoError(t, err)

	trMg := ds.NewTreeManager(&extractedRoot)
	it := ds.NewTreeIterator(trMg)
	cnt := 0
	for it.HasNext() {
		curNode, err := it.Next()
		got := curNode.Info
		require.NoError(t, err)
		_, ok := got.(*FileNode)
		require.True(t, ok)
		cnt++
	}
	require.Equal(t, 9, cnt)
}
