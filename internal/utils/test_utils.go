package utils

import (
	"fmt"
	"github.com/heroku/self/MetaManager/internal/ds"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func NewDirLifeCycleTester(t *testing.T, dir *MockDir, testExecFunc func(*testing.T, string)) DirLifeCycleTester {
	return DirLifeCycleTester{
		dir:          dir,
		t:            t,
		testExecFunc: testExecFunc,
	}
}

/*
Abstracts out the dir creation and destruction
while working on test which requires an actual directory
structure
*/
type DirLifeCycleTester struct {
	// dir contains directory structure relative to dirLoc
	dir *MockDir
	// dirLoc points to the absolute path of the parent of "dir"
	// creted after initializeDir
	parentDirLoc string
	t            *testing.T
	testExecFunc func(*testing.T, string)
}

func (d *DirLifeCycleTester) Execute() {
	err := d.initializeDir()
	require.NoError(d.t, err)
	defer d.deleteDir()

	d.testExecFunc(d.t, filepath.Join(d.parentDirLoc, d.dir.DirName))
}

func (d *DirLifeCycleTester) deleteDir() error {
	return os.RemoveAll(d.parentDirLoc)
}

func (d *DirLifeCycleTester) initializeDir() error {
	if d.dir == nil {
		return fmt.Errorf("uninitialized mock dir")
	}
	loc, err := CreateDirStructure(d.dir)
	if err != nil {
		return err
	}

	d.parentDirLoc = loc
	return nil
}

type MockDir struct {
	DirName string
	Files   []string
	Dirs    []*MockDir
}

/*
No need to save states, so keep it functional
Returns: <PathToDir>, <error>
*/
func CreateDirStructure(root *MockDir) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	randDirName := filepath.Join(wd, "testing_"+strconv.Itoa(rand.Int()))

	topLevelDir := MockDir{
		DirName: randDirName,
		Dirs:    []*MockDir{root},
	}

	queueDirs := []*MockDir{&topLevelDir}

	for i := 0; i < len(queueDirs); i++ {
		err := os.Mkdir(queueDirs[i].DirName, 0755)
		if err != nil {
			return "", err
		}

		for _, file := range queueDirs[i].Files {
			_, err = os.Create(filepath.Join(queueDirs[i].DirName, file))
			if err != nil {
				return "", err
			}
		}

		for _, dir := range queueDirs[i].Dirs {
			newDir := &MockDir{
				DirName: filepath.Join(queueDirs[i].DirName, dir.DirName),
				Dirs:    dir.Dirs,
				Files:   dir.Files,
			}
			queueDirs = append(queueDirs, newDir)
		}
	}

	return topLevelDir.DirName, nil
}

func ValidateNodeCnt(t *testing.T, node *ds.TreeNode, expCnt int) {
	if node == nil {
		require.Equal(t, expCnt, 0)
		return
	}
	it := ds.NewTreeIterator(ds.NewTreeManager(node))
	cnt := 0
	// fmt.Println("Start:")
	for it.HasNext() {
		// curNode, err := it.Next()
		_, err := it.Next()
		// val := curNode.Info
		require.NoError(t, err)

		// pr := val.(file.NodeInformable)

		// fmt.Println(pr.GetAbsPath())
		cnt += 1
	}
	// fmt.Println("End")
	require.Equal(t, expCnt, cnt)
}
