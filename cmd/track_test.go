package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/file"
	filesyspkg "github.com/heroku/self/MetaManager/internal/filesys"
	"github.com/heroku/self/MetaManager/internal/repository/filesys"
	"github.com/heroku/self/MetaManager/internal/services"
	servicemocks "github.com/heroku/self/MetaManager/internal/services/mocks"
	"github.com/heroku/self/MetaManager/internal/storage"
	"github.com/heroku/self/MetaManager/internal/utils"

	"github.com/stretchr/testify/mock"
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
			2, 3, 8, 8, 8, 16, // last: root + "*" (tracked nodes under root)
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

// setupMockGDriveService creates a mock GDrive service with test data and sets up expectations.
func setupMockGDriveService(t *testing.T) *servicemocks.MockGDriveServiceInterface {
	mockSvc := servicemocks.NewMockGDriveServiceInterface(t)

	// Set up test data structure:
	// root/
	//   Folder1/
	//     Sub/
	//       file3.txt
	//     file2.txt
	//   file1.txt

	// Root folder contents - will be called when tracking root
	mockSvc.On("ListFolder", mock.Anything, services.DriveRootID).
		Return([]services.RootEntry{
			{Id: "folder1", Name: "Folder1", IsFolder: true, MimeType: services.DriveFolderMimeType},
			{Id: "file1", Name: "file1.txt", IsFolder: false, MimeType: "text/plain"},
		}, nil).
		Maybe() // Maybe() allows this to be called 0 or more times

	// Folder1 contents - will be called when tracking /Folder1
	mockSvc.On("ListFolder", mock.Anything, "folder1").
		Return([]services.RootEntry{
			{Id: "sub1", Name: "Sub", IsFolder: true, MimeType: services.DriveFolderMimeType},
			{Id: "file2", Name: "file2.txt", IsFolder: false, MimeType: "text/plain"},
		}, nil).
		Maybe()

	// Sub folder contents - will be called when tracking recursively
	mockSvc.On("ListFolder", mock.Anything, "sub1").
		Return([]services.RootEntry{
			{Id: "file3", Name: "file3.txt", IsFolder: false, MimeType: "text/plain"},
		}, nil).
		Maybe()

	// ResolvePath expectations
	// Note: TrackGDrive trims slashes, so "/" becomes "", "/Folder1" becomes "Folder1", etc.
	mockSvc.On("ResolvePath", mock.Anything, "").
		Return(services.DriveRootID, nil).
		Maybe()

	mockSvc.On("ResolvePath", mock.Anything, "Folder1").
		Return("folder1", nil).
		Maybe()

	mockSvc.On("ResolvePath", mock.Anything, "Folder1/Sub").
		Return("sub1", nil).
		Maybe()

	return mockSvc
}

func TestTrackCmdGDrive(t *testing.T) {
	const gdriveCtxName = "gdrive-test"
	dir := t.TempDir()
	os.Setenv("MM_TEST_CONTEXT_DIR", dir)
	os.Setenv("MM_CONTEXT", gdriveCtxName)
	defer os.Unsetenv("MM_TEST_CONTEXT_DIR")
	defer os.Unsetenv("MM_CONTEXT")

	// Create gdrive context
	err := defaultStore.Create(gdriveCtxName, filesys.TypeGDrive)
	require.NoError(t, err)
	err = EnsureAppDataDir(gdriveCtxName)
	require.NoError(t, err)

	// Set up mock GDrive service with test data
	mockSvc := setupMockGDriveService(t)

	// Create scanner with mock service for path normalization tests
	scanner := filesyspkg.NewGDriveScanner(mockSvc)

	// Test path normalization
	path1, rec1 := scanner.NormalizeTrackPath("gdrive:/Folder1")
	require.Equal(t, "/Folder1", path1)
	require.False(t, rec1)

	path2, rec2 := scanner.NormalizeTrackPath("gdrive:/Folder1*")
	require.Equal(t, "/Folder1", path2)
	require.True(t, rec2)

	path3, rec3 := scanner.NormalizeTrackPath("/Folder1")
	require.Equal(t, "/Folder1", path3)
	require.False(t, rec3)

	path4, rec4 := scanner.NormalizeTrackPath("/Folder1/Sub*")
	require.Equal(t, "/Folder1/Sub", path4)
	require.True(t, rec4)

	path5, rec5 := scanner.NormalizeTrackPath("Folder1")
	require.Equal(t, "/Folder1", path5)
	require.False(t, rec5)

	// Create GDriveTracker directly with mock service for testing
	// This allows us to inject the mock service instead of relying on ContextAwareTracker
	// which would call the real services.GetGDriveService()
	tracker := filesyspkg.NewGDriveTracker(mockSvc)

	// Test tracking root (non-recursive)
	tree, err := tracker.Track("/")
	require.NoError(t, err)
	require.NotNil(t, tree)
	info, ok := tree.Info.(*file.FileNode)
	require.True(t, ok)
	require.Equal(t, file.GDrivePathPrefix, info.AbsPath)
	require.Len(t, tree.Children, 0) // Folder1 and file1

	// Test tracking a folder (non-recursive)
	tree2, err := tracker.Track("/Folder1")
	require.NoError(t, err)
	require.NotNil(t, tree2)
	info2, ok := tree2.Info.(*file.FileNode)
	require.True(t, ok)
	require.Equal(t, file.GDrivePathPrefix+"Folder1", info2.AbsPath)
	require.Len(t, tree2.Children, 0) // Sub and file2

	// Test tracking recursively
	tree3, err := tracker.Track("/Folder1*")
	require.NoError(t, err)
	require.NotNil(t, tree3)
	require.Len(t, tree3.Children, 2)
	// Find Sub folder by name
	var subNode *file.FileNode
	for _, child := range tree3.Children {
		if info := child.Info.(*file.FileNode); info.AbsPath == file.GDrivePathPrefix+"Folder1/Sub" {
			subNode = info
			require.Len(t, child.Children, 1) // file3 inside Sub
			break
		}
	}
	require.NotNil(t, subNode, "Sub folder should be found")

	// Test tracking root recursively
	tree4, err := tracker.Track("/*")
	require.NoError(t, err)
	require.NotNil(t, tree4)
	require.Len(t, tree4.Children, 2) // Folder1 and file1

	// Count total nodes in recursive root (should be: root + Folder1 + file1 + Sub + file2 + file3 = 6)
	nodeCount := countNodes(tree4)
	require.Equal(t, 6, nodeCount, "recursive root should have 6 nodes total")
}

// countNodes recursively counts all nodes in a tree.
func countNodes(node *ds.TreeNode) int {
	if node == nil {
		return 0
	}
	count := 1 // count self
	for _, child := range node.Children {
		count += countNodes(child)
	}
	return count
}
