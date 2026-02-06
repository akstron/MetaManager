package cmd

import (
	"github.com/heroku/self/MetaManager/internal/utils"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func InitializeRootAndScan(rootPath string) error {
	err := InitRoot(rootPath)
	if err != nil {
		return err
	}

	err = trackInternal(rootPath + "*")
	if err != nil {
		return err
	}

	return nil
}

func TestTagAddAndGetE2E(t *testing.T) {
	dirStructure := &utils.MockDir{
		DirName: "1_1",
		Files:   []string{"1_a", "1_b"},
		Dirs: []*utils.MockDir{
			{
				DirName: "2_1",
				Files:   []string{"2_a"},
			},
			{
				DirName: "2_2",
				Dirs: []*utils.MockDir{
					{
						DirName: "3_1",
					},
				},
			},
		},
	}

	testExecFunc := func(t *testing.T, root string) {
		os.Setenv("MM_TEST_ENV_DIR", root)
		err := InitializeRootAndScan(root)
		require.NoError(t, err)

		loc := filepath.Join(root, "1_a")
		loc2 := filepath.Join(root, "2_1", "2_a")
		loc3 := filepath.Join(root, "2_2")
		loc4 := filepath.Join(root, "2_2", "3_1")

		tags := []string{"hello", "world", "2", "random"}
		locs := []string{loc, loc2, loc3, loc4}
		for i, l := range locs {
			err := tagAddInternal([]string{l, tags[i]})
			require.NoError(t, err)
		}

		err = tagAddInternal([]string{loc3, "Hello World"})
		require.NoError(t, err)
		err = tagAddInternal([]string{loc, "Hello World"})
		require.NoError(t, err)

		for i, l := range locs {
			result, err := tagGetInternal(tags[i])
			require.NoError(t, err)
			require.Equal(t, 1, len(result))
			require.Equal(t, l, result[0])
		}

		result, err := tagGetInternal("Hello World")
		require.NoError(t, err)
		require.Equal(t, 2, len(result))
		sort.Strings(result)
		require.Equal(t, result[0], loc)
		require.Equal(t, result[1], loc3)
	}
	testExectutor := utils.NewDirLifeCycleTester(t, dirStructure, testExecFunc)
	testExectutor.Execute()
}
