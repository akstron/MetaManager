package ds

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type emptyStruct struct{}

func (*emptyStruct) Name() string {
	return "EMPTY"
}

func TestTreeNodeIteration(t *testing.T) {
	root := &TreeNode{
		Info: &emptyStruct{},
		Children: []*TreeNode{
			{
				Info: &emptyStruct{},
			},
			{
				Info: &emptyStruct{},
				Children: []*TreeNode{
					{
						Info: &emptyStruct{},
						Children: []*TreeNode{
							{
								Info: &emptyStruct{},
							},
						},
					},
					{
						Info: &emptyStruct{},
					},
					{
						Info: &emptyStruct{},
					},
					{
						Info: &emptyStruct{},
					},
				},
			},
			{
				Info: &emptyStruct{},
			},
		},
	}

	trMg := NewTreeManager(root)
	it := NewTreeIterator(trMg)
	cnt := 0

	for it.HasNext() {
		_, err := it.Next()
		require.NoError(t, err)
		cnt += 1
	}

	require.Equal(t, 9, cnt)
}

func TestTreeNodeIterationWithNilEntries(t *testing.T) {

	root := &TreeNode{
		Info: &emptyStruct{},
		Children: []*TreeNode{
			{
				Info: &emptyStruct{},
			},
			{
				Info: &emptyStruct{},
				Children: []*TreeNode{
					{
						Info: &emptyStruct{},
						Children: []*TreeNode{
							{
								Info: &emptyStruct{},
							},
						},
					},
					{
						Info: nil,
					},
					{
						Info: &emptyStruct{},
					},
					{
						Info: nil,
					},
				},
			},
			{
				Info: &emptyStruct{},
			},
		},
	}

	trMg := NewTreeManager(root)
	it := NewTreeIterator(trMg)
	cnt := 0

	for it.HasNext() {
		_, err := it.Next()
		require.NoError(t, err)
		cnt += 1
	}

	require.Equal(t, 7, cnt)
}
