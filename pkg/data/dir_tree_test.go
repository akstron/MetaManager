package data

import (
	"github/akstron/MetaManager/ds"
	"github/akstron/MetaManager/pkg/utils"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMergeNodeWithPath(t *testing.T) {
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

		locs := []string{"1_1", "1_1/1_a", "1_1", "1_1/2_2/3_1/temp", "1_1/2_2/x", "1_1/2_1/2_a", "1_1/2_2/3_1"}
		outputs := []int{1, 2, 2, 5, 6, 8, 8}

		dm := NewDirTreeManager(ds.NewTreeManager(nil))

		for i, loc := range locs {
			absLoc := filepath.Join(root, "..", loc)
			err := dm.MergeNodeWithPath(absLoc)
			require.NoError(t, err)

			utils.ValidateNodeCnt(t, dm.Root, outputs[i])
		}
	}
	testExectutor := utils.NewDirLifeCycleTester(t, dirStructure, testExecFunc)
	testExectutor.Execute()
}
