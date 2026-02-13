package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/heroku/self/MetaManager/internal/repository/filesys"
	"github.com/heroku/self/MetaManager/internal/storage"
	"github.com/heroku/self/MetaManager/internal/utils"

	"github.com/stretchr/testify/require"
)

func TestTrackCmd(t *testing.T) {
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
						Dirs: []*utils.MockDir{
							{
								DirName: "temp",
							},
						},
					},
				},
				Files: []string{"x", "y", "z"},
			},
		},
	}

	testExecFunc := func(t *testing.T, root string) {
		os.Setenv("MM_TEST_CONTEXT_DIR", root)
		os.Setenv("MM_CONTEXT", "default")
		defer os.Unsetenv("MM_TEST_CONTEXT_DIR")
		defer os.Unsetenv("MM_CONTEXT")

		err := defaultStore.Create("default", filesys.TypeLocal)
		require.NoError(t, err)
		err = EnsureAppDataDir("default")
		require.NoError(t, err)

		rw, err := storage.GetRW("default")
		require.NoError(t, err)

		loc := filepath.Join(root, "1_a")
		loc2 := filepath.Join(root, "2_2")
		loc3 := filepath.Join(root, "2_2")
		loc3 = loc3 + "*"

		locs := []string{
			loc, loc2, loc3,
			filepath.Join(root, "2_2", "x"),
			filepath.Join(root, "2_2", "3_1") + "*",
			root + "*",
		}

		outputs := []int{
			2, 3, 8, 8, 8, 16, /*includes .mm/default created for context*/
		}

		for i, loc := range locs {
			err = trackInternal("default", loc)
			require.NoError(t, err)

			node, err := rw.Read()
			require.NoError(t, err)

			utils.ValidateNodeCnt(t, node, outputs[i])
		}
	}
	testExectutor := utils.NewDirLifeCycleTester(t, dirStructure, testExecFunc)
	testExectutor.Execute()
}
