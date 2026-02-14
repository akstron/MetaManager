package cmd

import (
	"os"
	"testing"

	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/file"
	"github.com/heroku/self/MetaManager/internal/repository/filesys"
	"github.com/heroku/self/MetaManager/internal/storage"
	"github.com/heroku/self/MetaManager/internal/utils"

	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
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
		os.Setenv("MM_CONTEXT", "default")
		defer os.Unsetenv("MM_TEST_CONTEXT_DIR")
		defer os.Unsetenv("MM_CONTEXT")
		err := EnsureAppDataDir("default")
		require.NoError(t, err)

		rw, err := storage.GetRW("default")
		require.NoError(t, err)

		node, err := rw.Read()
		require.NoError(t, err)

		it := ds.NewTreeIterator(ds.NewTreeManager(node))
		cnt := 0
		var got *file.FileNode

		for it.HasNext() {
			curNode, err := it.Next()
			info := curNode.Info
			require.NoError(t, err)
			if fn, ok := info.(*file.FileNode); !ok {
				t.FailNow()
			} else {
				got = fn
			}
			cnt += 1
		}
		require.Equal(t, 1, cnt)
		// When MM_TEST_CONTEXT_DIR is set, root is set to that path so merge in track tests works.
		require.Equal(t, got.AbsPath, root)
	}
	testExectutor := utils.NewDirLifeCycleTester(t, dirStructure, testExecFunc)
	testExectutor.Execute()
}

// TestEnsureAppDataDirGDriveRoot asserts that when context type is gdrive, the initial root AbsPath is gdrive:/
func TestEnsureAppDataDirGDriveRoot(t *testing.T) {
	const gdriveCtxName = "gdrive-ctx"
	dir := t.TempDir()
	os.Setenv("MM_TEST_CONTEXT_DIR", dir)
	os.Setenv("MM_CONTEXT", gdriveCtxName)
	defer os.Unsetenv("MM_TEST_CONTEXT_DIR")
	defer os.Unsetenv("MM_CONTEXT")

	err := defaultStore.Create(gdriveCtxName, filesys.TypeGDrive)
	require.NoError(t, err)
	err = EnsureAppDataDir(gdriveCtxName)
	require.NoError(t, err)

	rw, err := storage.GetRW(gdriveCtxName)
	require.NoError(t, err)
	node, err := rw.Read()
	require.NoError(t, err)
	require.NotNil(t, node)

	info, ok := node.Info.(*file.FileNode)
	require.True(t, ok, "root info should be *file.FileNode")
	require.Equal(t, file.GDrivePathPrefix, info.AbsPath, "gdrive context root AbsPath should be gdrive:/")
}
