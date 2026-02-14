package filesys

import (
	"os"
	"path/filepath"
	"testing"

	contextrepo "github.com/heroku/self/MetaManager/internal/repository/filesys"
	"github.com/heroku/self/MetaManager/internal/utils"

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
		os.Setenv("MM_TEST_CONTEXT_DIR", root)
		defer os.Unsetenv("MM_TEST_CONTEXT_DIR")

		// Set up context repository and create a default local context
		ctxRepo := contextrepo.NewContextRepositoryImpl(nil)
		err := ctxRepo.Create("default", contextrepo.TypeLocal)
		require.NoError(t, err)
		err = ctxRepo.SetCurrent("default")
		require.NoError(t, err)

		// Create tracker
		tracker := NewContextAwareTracker(ctxRepo)

		loc := filepath.Join(root, "1_a")
		loc2 := filepath.Join(root, "2_2")
		loc3 := filepath.Join(root, "2_2")
		loc3 = loc3 + "*"

		node, err := tracker.Track(loc)
		require.NoError(t, err)
		utils.ValidateNodeCnt(t, node, 1)

		node, err = tracker.Track(loc2)
		require.NoError(t, err)
		utils.ValidateNodeCnt(t, node, 1)

		node, err = tracker.Track(loc3)
		require.NoError(t, err)
		utils.ValidateNodeCnt(t, node, 6)
	}
	testExectutor := utils.NewDirLifeCycleTester(t, dirStructure, testExecFunc)
	testExectutor.Execute()
}
