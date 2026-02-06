package filesys

import (
	"github.com/heroku/self/MetaManager/internal/utils"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrack(t *testing.T) {
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

		loc := filepath.Join(root, "1_a")
		loc2 := filepath.Join(root, "2_2")
		loc3 := filepath.Join(root, "2_2")
		loc3 = loc3 + "*"

		node, err := Track(loc)
		require.NoError(t, err)
		utils.ValidateNodeCnt(t, node, 1)

		node, err = Track(loc2)
		require.NoError(t, err)
		utils.ValidateNodeCnt(t, node, 1)

		node, err = Track(loc3)
		require.NoError(t, err)
		utils.ValidateNodeCnt(t, node, 6)
	}
	testExectutor := utils.NewDirLifeCycleTester(t, dirStructure, testExecFunc)
	testExectutor.Execute()
}
