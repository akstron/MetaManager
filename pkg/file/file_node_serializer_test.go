package file

import (
	"encoding/json"
	"github/akstron/MetaManager/ds"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDataSerialization(t *testing.T) {
	root := &ds.TreeNode{
		Info: &DirNode{},
		Children: []*ds.TreeNode{
			{
				Info: &FileNode{},
			},
			{
				Info: &DirNode{},
				Children: []*ds.TreeNode{
					{
						Info: &DirNode{},
						Children: []*ds.TreeNode{
							{
								Info: &FileNode{},
							},
						},
					},
					{
						Info: &DirNode{},
					},
					{
						Info: &DirNode{},
					},
					{
						Info: &FileNode{},
					},
				},
			},
			{
				Info: &FileNode{},
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
	dirCnt := 0
	fileCnt := 0

	for it.HasNext() {
		got, err := it.Next()
		require.NoError(t, err)

		switch got.(type) {
		case *DirNode:
			dirCnt += 1
		case *FileNode:
			fileCnt += 1
		default:
			t.FailNow()
		}
	}

	require.Equal(t, 5, dirCnt)
	require.Equal(t, 4, fileCnt)
}
