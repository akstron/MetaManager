package cmd

import (
	"github.com/heroku/self/MetaManager/internal/utils"
	"github.com/heroku/self/MetaManager/internal/storage"
	"os"
	"path/filepath"
	"testing"

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
		os.Setenv("MM_TEST_ENV_DIR", root)
		err := InitRoot(root)
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
			2, 3, 8, 8, 8, 15, /*due to .mm creation*/
		}

		rw, err := storage.GetRW()
		require.NoError(t, err)

		for i, loc := range locs {
			err = trackInternal(loc)
			require.NoError(t, err)

			node, err := rw.Read()
			require.NoError(t, err)

			utils.ValidateNodeCnt(t, node, outputs[i])
		}
	}
	testExectutor := utils.NewDirLifeCycleTester(t, dirStructure, testExecFunc)
	testExectutor.Execute()
}
