package cmd

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/heroku/self/MetaManager/internal/repository/filesys"
	"github.com/heroku/self/MetaManager/internal/utils"

	"github.com/stretchr/testify/require"
)

func InitializeRootAndScan(rootPath string) error {
	os.Setenv("MM_TEST_CONTEXT_DIR", rootPath)
	os.Setenv("MM_CONTEXT", "default")
	if err := defaultStore.Create("default", filesys.TypeLocal); err != nil {
		return err
	}
	if err := EnsureAppDataDir("default"); err != nil {
		return err
	}
	return trackInternal("default", rootPath+"*")
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
		os.Setenv("MM_TEST_CONTEXT_DIR", root)
		defer os.Unsetenv("MM_TEST_CONTEXT_DIR")
		err := InitializeRootAndScan(root)
		require.NoError(t, err)

		loc := filepath.Join(root, "1_a")
		loc2 := filepath.Join(root, "2_1", "2_a")
		loc3 := filepath.Join(root, "2_2")
		loc4 := filepath.Join(root, "2_2", "3_1")

		tags := []string{"hello", "world", "2", "random"}
		locs := []string{loc, loc2, loc3, loc4}
		for i, l := range locs {
			err := tagAddInternal("default", []string{l, tags[i]})
			require.NoError(t, err)
		}

		err = tagAddInternal("default", []string{loc3, "Hello World"})
		require.NoError(t, err)
		err = tagAddInternal("default", []string{loc, "Hello World"})
		require.NoError(t, err)

		for i, l := range locs {
			result, err := tagSearchInternal("default", tags[i])
			require.NoError(t, err)
			require.Equal(t, 1, len(result))
			require.Equal(t, l, result[0])
		}

		result, err := tagSearchInternal("default", "Hello World")
		require.NoError(t, err)
		require.Equal(t, 2, len(result))
		sort.Strings(result)
		require.Equal(t, result[0], loc)
		require.Equal(t, result[1], loc3)

		getResult, err := tagGetInternal("default", loc)
		require.NoError(t, err)
		require.Equal(t, 2, len(getResult))
		sort.Strings(getResult)
		require.Equal(t, getResult[0], "Hello World")
		require.Equal(t, getResult[1], tags[0])

		getResult, err = tagGetInternal("default", loc3)
		require.NoError(t, err)
		require.Equal(t, 2, len(getResult))
		sort.Strings(getResult)
		require.Equal(t, getResult[1], "Hello World")
		require.Equal(t, getResult[0], tags[2])

		getResult, err = tagGetInternal("default", loc2)
		require.NoError(t, err)
		require.Equal(t, 1, len(getResult))
		require.Equal(t, getResult[0], tags[1])

		getResult, err = tagGetInternal("default", loc4)
		require.NoError(t, err)
		require.Equal(t, 1, len(getResult))
		require.Equal(t, getResult[0], tags[3])
	}
	testExectutor := utils.NewDirLifeCycleTester(t, dirStructure, testExecFunc)
	testExectutor.Execute()
}
