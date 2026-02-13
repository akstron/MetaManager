package data

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/utils"
	"github.com/stretchr/testify/require"
)

const testContextEnvVar = "MM_CONTEXT"

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

func TestSplitNodeWithPath(t *testing.T) {
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

		splitLocs := []string{"1_1/2_2/3_1", "1_1/2_2", "1_1/1_a", "1_1"}
		splitOutputs := []int{6, 4, 3, 0}

		for i, loc := range splitLocs {
			absLoc := filepath.Join(root, "..", loc)
			err := dm.SplitNodeWithPath(absLoc)
			require.NoError(t, err)

			utils.ValidateNodeCnt(t, dm.Root, splitOutputs[i])
		}
	}
	testExectutor := utils.NewDirLifeCycleTester(t, dirStructure, testExecFunc)
	testExectutor.Execute()

}
