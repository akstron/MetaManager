package ds

// var _ = (json.Marshaler(&TreeNode{}))

import (
	"encoding/json"
	"testing"

	"fmt"

	"github.com/stretchr/testify/require"
)

type MockNodeInfo struct {
	MName string
}

func (m *MockNodeInfo) Name() string {
	return m.MName
}

func TestTreeNodeSerialization(t *testing.T) {
	root := &TreeNode{
		Info: &MockNodeInfo{MName: "root"},
		Children: []*TreeNode{
			{
				Info: &MockNodeInfo{MName: "child1"},
			},
			{
				Info: &MockNodeInfo{MName: "child2"},
			},
		},
	}

	serializedNode, err := json.Marshal(root)
	require.NoError(t, err)

	fmt.Println(string(serializedNode))

	var extractedRoot TreeNodeJSON
	err = json.Unmarshal(serializedNode, &extractedRoot)
	require.NoError(t, err)

	require.Equal(t, root.Info.Name(), extractedRoot.Info["MName"])
	require.Equal(t, len(root.Children), len(extractedRoot.Children))
}
