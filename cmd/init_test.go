package cmd

import (
	"github/akstron/MetaManager/ds"
	"github/akstron/MetaManager/pkg/file"
	"github/akstron/MetaManager/pkg/utils"
	"github/akstron/MetaManager/storage"
	"os"
	"testing"

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
		os.Setenv("MM_TEST_ENV_DIR", root)
		err := InitRoot(root)
		require.NoError(t, err)

		rw, err := storage.GetRW()
		require.NoError(t, err)

		node, err := rw.Read()
		require.NoError(t, err)

		it := ds.NewTreeIterator(ds.NewTreeManager(node))
		cnt := 0
		var got *file.DirNode

		for it.HasNext() {
			info, err := it.Next()
			require.NoError(t, err)
			if dirNode, ok := info.(*file.DirNode); !ok {
				t.FailNow()
			} else {
				got = dirNode
			}
			cnt += 1
		}
		require.Equal(t, 1, cnt)
		require.Equal(t, got.AbsPath, root)
	}
	testExectutor := utils.NewDirLifeCycleTester(t, dirStructure, testExecFunc)
	testExectutor.Execute()
}
