package data

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/file"
	"github.com/heroku/self/MetaManager/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestUnixPathSplitterImpl_Split(t *testing.T) {
	splitter := &UnixPathSplitterImpl{}

	_, err := splitter.Split("")
	require.Error(t, err)
	require.Contains(t, err.Error(), "path is not a unix path")

	segments, err := splitter.Split("/")
	require.NoError(t, err)
	require.Empty(t, segments)

	_, err = splitter.Split("a")
	require.Error(t, err)
	require.Contains(t, err.Error(), "path is not a unix path")

	_, err = splitter.Split("/a/")
	require.Error(t, err)
	require.Contains(t, err.Error(), "path cannot end with a slash")

	_, err = splitter.Split("/a/b/")
	require.Error(t, err)
	require.Contains(t, err.Error(), "path cannot end with a slash")

	segments, err = splitter.Split("/a")
	require.NoError(t, err)
	require.Equal(t, []string{"a"}, segments)

	segments, err = splitter.Split("/a/b/c")
	require.NoError(t, err)
	require.Equal(t, []string{"a", "b", "c"}, segments)

	_, err = splitter.Split("a/b")
	require.Error(t, err)
	require.Contains(t, err.Error(), "path is not a unix path")
}

func TestGDrivePathSplitterImpl_Split(t *testing.T) {
	splitter := &GDrivePathSplitterImpl{}

	_, err := splitter.Split("/not/gdrive")
	require.Error(t, err)
	require.Contains(t, err.Error(), "path is not a gdrive path")

	_, err = splitter.Split(file.GDrivePathPrefix)
	require.Error(t, err)
	require.Contains(t, err.Error(), "path contains empty segments")

	segments, err := splitter.Split(file.GDrivePathPrefix + "Folder")
	require.NoError(t, err)
	require.Equal(t, []string{"gdrive:", "Folder"}, segments)

	segments, err = splitter.Split(file.GDrivePathPrefix + "Folder/Sub/file.pdf")
	require.NoError(t, err)
	require.Equal(t, []string{"gdrive:", "Folder", "Sub", "file.pdf"}, segments)

	_, err = splitter.Split(file.GDrivePathPrefix + "/")
	require.Error(t, err)
	require.Contains(t, err.Error(), "path contains empty segments")

	_, err = splitter.Split(file.GDrivePathPrefix + "Folder/")
	require.Error(t, err)
	require.Contains(t, err.Error(), "path contains empty segments")

	_, err = splitter.Split(file.GDrivePathPrefix + "Folder/Sub/")
	require.Error(t, err)
	require.Contains(t, err.Error(), "path contains empty segments")
}

const testContextEnvVar = "MM_CONTEXT"

func TestGetPathSplitter(t *testing.T) {
	// Restore env after test so we don't affect other tests.
	prev := os.Getenv(testContextEnvVar)
	defer func() {
		if prev == "" {
			os.Unsetenv(testContextEnvVar)
		} else {
			os.Setenv(testContextEnvVar, prev)
		}
	}()

	os.Setenv(testContextEnvVar, "local")
	splitter, err := GetPathSplitter()
	require.NoError(t, err)
	require.IsType(t, &UnixPathSplitterImpl{}, splitter)
	segments, err := splitter.Split("/local/path")
	require.NoError(t, err)
	require.Equal(t, []string{"local", "path"}, segments)

	os.Setenv(testContextEnvVar, "gdrive")
	splitter, err = GetPathSplitter()
	require.NoError(t, err)
	require.IsType(t, &GDrivePathSplitterImpl{}, splitter)
	segments, err = splitter.Split(file.GDrivePathPrefix + "My Drive/doc")
	require.NoError(t, err)
	require.Equal(t, []string{"gdrive:", "My Drive", "doc"}, segments)
}

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
