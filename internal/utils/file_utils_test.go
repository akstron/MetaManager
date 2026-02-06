package utils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindMMDirPath(t *testing.T) {
	dirStructure := &MockDir{
		DirName: "1_1",
		Files:   []string{"1_a", "1_b"},
		Dirs: []*MockDir{
			{
				DirName: "2_1",
				Files:   []string{"2_a"},
			},
			{
				DirName: "2_2",
				Dirs: []*MockDir{
					{
						DirName: "3_1",
					},
				},
			},
			{
				DirName: ".mm",
			},
		},
	}

	testExecFunc := func(t *testing.T, root string) {
		curWD := filepath.Join(root, "2_2", "3_1")
		found, path, err := findMMDirPathInternal(curWD)
		require.NoError(t, err)
		require.True(t, found)
		require.Equal(t, path, filepath.Join(root, ".mm"))
	}
	testExectutor := NewDirLifeCycleTester(t, dirStructure, testExecFunc)
	testExectutor.Execute()
}
