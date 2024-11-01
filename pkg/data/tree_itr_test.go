package data

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTreeNodeIteration(t *testing.T) {
	root := &DirNode{
		FileChildren: []*FileNode{
			{}, {}, {},
		},
		DirChildren: []*DirNode{
			{},
			{
				DirChildren: []*DirNode{
					{}, {},
				},
				FileChildren: []*FileNode{
					{},
				},
			},
		},
	}

	trMg := NewTreeManager()
	trMg.Root = root

	it := NewTreeIterator(&trMg)
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

func TestTreeNodeIterationWithNilEntries(t *testing.T) {
	root := &DirNode{
		FileChildren: []*FileNode{
			{}, nil, {},
		},
		DirChildren: []*DirNode{
			{},
			{
				DirChildren: []*DirNode{
					nil, {},
				},
				FileChildren: []*FileNode{
					{},
				},
			},
		},
	}

	trMg := NewTreeManager()
	trMg.Root = root

	it := NewTreeIterator(&trMg)
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

	require.Equal(t, 4, dirCnt)
	require.Equal(t, 3, fileCnt)
}
