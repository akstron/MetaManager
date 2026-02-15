package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/file"
)

func TestFileStorageRW_Write(t *testing.T) {
	t.Run("write simple tree with root node", func(t *testing.T) {
		dir := t.TempDir()
		dataFilePath := filepath.Join(dir, "data.json")

		rw, err := NewFileStorageRW(dataFilePath)
		require.NoError(t, err)

		root := &ds.TreeNode{
			Info: &file.FileNode{
				GeneralNode: file.GeneralNode{
					AbsPath: "/root",
					Tags:    []string{"tag1"},
					Id:      "root-id",
				},
			},
			Children: nil,
		}

		err = rw.Write(root)
		require.NoError(t, err)

		// Verify file was created
		_, err = os.Stat(dataFilePath)
		require.NoError(t, err)

		// Verify file contents
		data, err := os.ReadFile(dataFilePath)
		require.NoError(t, err)

		var treeNodeJSON ds.TreeNodeJSON
		err = json.Unmarshal(data, &treeNodeJSON)
		require.NoError(t, err)

		require.NotNil(t, treeNodeJSON.Info)
		absPath, ok := treeNodeJSON.Info["AbsPath"]
		require.True(t, ok, "AbsPath key should exist in Info map")
		require.Equal(t, "/root", absPath)
		id, ok := treeNodeJSON.Info["Id"]
		require.True(t, ok, "Id key should exist in Info map")
		require.Equal(t, "root-id", id)
		require.Empty(t, treeNodeJSON.Children)
	})

	t.Run("write tree with children", func(t *testing.T) {
		dir := t.TempDir()
		dataFilePath := filepath.Join(dir, "data.json")

		rw, err := NewFileStorageRW(dataFilePath)
		require.NoError(t, err)

		root := &ds.TreeNode{
			Info: &file.FileNode{
				GeneralNode: file.GeneralNode{
					AbsPath: "/root",
				},
			},
			Children: []*ds.TreeNode{
				{
					Info: &file.FileNode{
						GeneralNode: file.GeneralNode{
							AbsPath: "/root/child1",
							Tags:    []string{"child-tag"},
							Id:      "child1-id",
						},
					},
					Children: nil,
				},
				{
					Info: &file.FileNode{
						GeneralNode: file.GeneralNode{
							AbsPath: "/root/child2",
						},
					},
					Children: []*ds.TreeNode{
						{
							Info: &file.FileNode{
								GeneralNode: file.GeneralNode{
									AbsPath: "/root/child2/grandchild",
									Id:      "grandchild-id",
								},
							},
							Children: nil,
						},
					},
				},
			},
		}

		err = rw.Write(root)
		require.NoError(t, err)

		// Verify file contents
		data, err := os.ReadFile(dataFilePath)
		require.NoError(t, err)

		var treeNodeJSON ds.TreeNodeJSON
		err = json.Unmarshal(data, &treeNodeJSON)
		require.NoError(t, err)

		absPath, ok := treeNodeJSON.Info["AbsPath"]
		require.True(t, ok)
		require.Equal(t, "/root", absPath)
		require.Len(t, treeNodeJSON.Children, 2)
		child1AbsPath, ok := treeNodeJSON.Children[0].Info["AbsPath"]
		require.True(t, ok)
		require.Equal(t, "/root/child1", child1AbsPath)
		child2AbsPath, ok := treeNodeJSON.Children[1].Info["AbsPath"]
		require.True(t, ok)
		require.Equal(t, "/root/child2", child2AbsPath)
		require.Len(t, treeNodeJSON.Children[1].Children, 1)
		grandchildAbsPath, ok := treeNodeJSON.Children[1].Children[0].Info["AbsPath"]
		require.True(t, ok)
		require.Equal(t, "/root/child2/grandchild", grandchildAbsPath)
	})

	t.Run("write Google Drive node", func(t *testing.T) {
		dir := t.TempDir()
		dataFilePath := filepath.Join(dir, "data.json")

		rw, err := NewFileStorageRW(dataFilePath)
		require.NoError(t, err)

		root := &ds.TreeNode{
			Info: &file.FileNode{
				GeneralNode: file.GeneralNode{
					AbsPath: file.GDrivePathPrefix + "Folder/file.txt",
					Tags:    []string{"gdrive"},
					Id:      "drive-id",
				},
				DriveId: "1a2b3c4d5e6f",
			},
			Children: nil,
		}

		err = rw.Write(root)
		require.NoError(t, err)

		// Verify file contents
		data, err := os.ReadFile(dataFilePath)
		require.NoError(t, err)

		var treeNodeJSON ds.TreeNodeJSON
		err = json.Unmarshal(data, &treeNodeJSON)
		require.NoError(t, err)

		absPath, ok := treeNodeJSON.Info["AbsPath"]
		require.True(t, ok)
		require.Equal(t, file.GDrivePathPrefix+"Folder/file.txt", absPath)
		driveId, ok := treeNodeJSON.Info["DriveId"]
		require.True(t, ok)
		require.Equal(t, "1a2b3c4d5e6f", driveId)
		id, ok := treeNodeJSON.Info["Id"]
		require.True(t, ok)
		require.Equal(t, "drive-id", id)
	})

	t.Run("write empty tree", func(t *testing.T) {
		dir := t.TempDir()
		dataFilePath := filepath.Join(dir, "data.json")

		rw, err := NewFileStorageRW(dataFilePath)
		require.NoError(t, err)

		root := &ds.TreeNode{
			Info: &file.FileNode{
				GeneralNode: file.GeneralNode{
					AbsPath: "/",
				},
			},
			Children: []*ds.TreeNode{},
		}

		err = rw.Write(root)
		require.NoError(t, err)

		// Verify file was created
		_, err = os.Stat(dataFilePath)
		require.NoError(t, err)
	})
}

func TestFileStorageRW_Read(t *testing.T) {
	t.Run("read simple tree with root node", func(t *testing.T) {
		dir := t.TempDir()
		dataFilePath := filepath.Join(dir, "data.json")

		// Create JSON file manually
		treeNodeJSON := ds.TreeNodeJSON{
			Info: map[string]interface{}{
				"AbsPath": "/root",
				"Tags":    []interface{}{"tag1"},
				"Id":      "root-id",
				"DriveId": "",
			},
			Children: []*ds.TreeNodeJSON{},
		}

		data, err := json.Marshal(treeNodeJSON)
		require.NoError(t, err)
		err = os.WriteFile(dataFilePath, data, 0666)
		require.NoError(t, err)

		rw, err := NewFileStorageRW(dataFilePath)
		require.NoError(t, err)

		node, err := rw.Read()
		require.NoError(t, err)
		require.NotNil(t, node)

		fn, ok := node.Info.(*file.FileNode)
		require.True(t, ok, "Info should be *file.FileNode")
		require.Equal(t, "/root", fn.AbsPath)
		require.Equal(t, []string{"tag1"}, fn.Tags)
		require.Equal(t, "root-id", fn.Id)
		require.Empty(t, fn.DriveId)
		require.Empty(t, node.Children)
	})

	t.Run("read tree with children", func(t *testing.T) {
		dir := t.TempDir()
		dataFilePath := filepath.Join(dir, "data.json")

		// Create JSON file manually
		treeNodeJSON := ds.TreeNodeJSON{
			Info: map[string]interface{}{
				"AbsPath": "/root",
			},
			Children: []*ds.TreeNodeJSON{
				{
					Info: map[string]interface{}{
						"AbsPath": "/root/child1",
						"Tags":    []interface{}{"child-tag"},
						"Id":      "child1-id",
						"DriveId": "",
					},
					Children: []*ds.TreeNodeJSON{},
				},
				{
					Info: map[string]interface{}{
						"AbsPath": "/root/child2",
						"DriveId": "",
					},
					Children: []*ds.TreeNodeJSON{
						{
							Info: map[string]interface{}{
								"AbsPath": "/root/child2/grandchild",
								"Id":      "grandchild-id",
								"DriveId": "",
							},
							Children: []*ds.TreeNodeJSON{},
						},
					},
				},
			},
		}

		data, err := json.Marshal(treeNodeJSON)
		require.NoError(t, err)
		err = os.WriteFile(dataFilePath, data, 0666)
		require.NoError(t, err)

		rw, err := NewFileStorageRW(dataFilePath)
		require.NoError(t, err)

		node, err := rw.Read()
		require.NoError(t, err)
		require.NotNil(t, node)

		fn, ok := node.Info.(*file.FileNode)
		require.True(t, ok)
		require.Equal(t, "/root", fn.AbsPath)
		require.Len(t, node.Children, 2)

		child1 := node.Children[0]
		child1Fn, ok := child1.Info.(*file.FileNode)
		require.True(t, ok)
		require.Equal(t, "/root/child1", child1Fn.AbsPath)
		require.Equal(t, []string{"child-tag"}, child1Fn.Tags)
		require.Equal(t, "child1-id", child1Fn.Id)

		child2 := node.Children[1]
		child2Fn, ok := child2.Info.(*file.FileNode)
		require.True(t, ok)
		require.Equal(t, "/root/child2", child2Fn.AbsPath)
		require.Len(t, child2.Children, 1)

		grandchild := child2.Children[0]
		grandchildFn, ok := grandchild.Info.(*file.FileNode)
		require.True(t, ok)
		require.Equal(t, "/root/child2/grandchild", grandchildFn.AbsPath)
		require.Equal(t, "grandchild-id", grandchildFn.Id)
	})

	t.Run("read Google Drive node", func(t *testing.T) {
		dir := t.TempDir()
		dataFilePath := filepath.Join(dir, "data.json")

		treeNodeJSON := ds.TreeNodeJSON{
			Info: map[string]interface{}{
				"AbsPath": file.GDrivePathPrefix + "Folder/file.txt",
				"Tags":    []interface{}{"gdrive"},
				"Id":      "drive-id",
				"DriveId": "1a2b3c4d5e6f",
			},
			Children: []*ds.TreeNodeJSON{},
		}

		data, err := json.Marshal(treeNodeJSON)
		require.NoError(t, err)
		err = os.WriteFile(dataFilePath, data, 0666)
		require.NoError(t, err)

		rw, err := NewFileStorageRW(dataFilePath)
		require.NoError(t, err)

		node, err := rw.Read()
		require.NoError(t, err)
		require.NotNil(t, node)

		fn, ok := node.Info.(*file.FileNode)
		require.True(t, ok)
		require.Equal(t, file.GDrivePathPrefix+"Folder/file.txt", fn.AbsPath)
		require.Equal(t, "1a2b3c4d5e6f", fn.DriveId)
		require.Equal(t, "drive-id", fn.Id)
		require.Equal(t, []string{"gdrive"}, fn.Tags)
	})

	t.Run("read non-existent file", func(t *testing.T) {
		dir := t.TempDir()
		dataFilePath := filepath.Join(dir, "nonexistent.json")

		rw, err := NewFileStorageRW(dataFilePath)
		require.NoError(t, err)

		node, err := rw.Read()
		require.Error(t, err)
		require.Nil(t, node)
		require.Contains(t, err.Error(), "no such file")
	})

	t.Run("read invalid JSON", func(t *testing.T) {
		dir := t.TempDir()
		dataFilePath := filepath.Join(dir, "data.json")

		invalidJSON := []byte(`{invalid json}`)
		err := os.WriteFile(dataFilePath, invalidJSON, 0666)
		require.NoError(t, err)

		rw, err := NewFileStorageRW(dataFilePath)
		require.NoError(t, err)

		node, err := rw.Read()
		require.Error(t, err)
		require.Nil(t, node)
	})

	t.Run("read empty file", func(t *testing.T) {
		dir := t.TempDir()
		dataFilePath := filepath.Join(dir, "data.json")

		err := os.WriteFile(dataFilePath, []byte{}, 0666)
		require.NoError(t, err)

		rw, err := NewFileStorageRW(dataFilePath)
		require.NoError(t, err)

		node, err := rw.Read()
		require.Error(t, err)
		require.Nil(t, node)
	})
}

func TestFileStorageRW_RoundTrip(t *testing.T) {
	t.Run("write and read simple tree", func(t *testing.T) {
		dir := t.TempDir()
		dataFilePath := filepath.Join(dir, "data.json")

		rw, err := NewFileStorageRW(dataFilePath)
		require.NoError(t, err)

		original := &ds.TreeNode{
			Info: &file.FileNode{
				GeneralNode: file.GeneralNode{
					AbsPath: "/test/path",
					Tags:    []string{"tag1", "tag2"},
					Id:      "test-id",
				},
				DriveId: "",
			},
			Children: nil,
		}

		err = rw.Write(original)
		require.NoError(t, err)

		read, err := rw.Read()
		require.NoError(t, err)
		require.NotNil(t, read)

		originalFn := original.Info.(*file.FileNode)
		readFn := read.Info.(*file.FileNode)

		require.Equal(t, originalFn.AbsPath, readFn.AbsPath)
		require.Equal(t, originalFn.Tags, readFn.Tags)
		require.Equal(t, originalFn.Id, readFn.Id)
		require.Equal(t, originalFn.DriveId, readFn.DriveId)
		require.Equal(t, len(original.Children), len(read.Children))
	})

	t.Run("write and read complex tree with nested children", func(t *testing.T) {
		dir := t.TempDir()
		dataFilePath := filepath.Join(dir, "data.json")

		rw, err := NewFileStorageRW(dataFilePath)
		require.NoError(t, err)

		original := &ds.TreeNode{
			Info: &file.FileNode{
				GeneralNode: file.GeneralNode{
					AbsPath: "/root",
					Tags:    []string{"root-tag"},
					Id:      "root-id",
				},
			},
			Children: []*ds.TreeNode{
				{
					Info: &file.FileNode{
						GeneralNode: file.GeneralNode{
							AbsPath: "/root/child1",
							Tags:    []string{"child1-tag"},
							Id:      "child1-id",
						},
					},
					Children: nil,
				},
				{
					Info: &file.FileNode{
						GeneralNode: file.GeneralNode{
							AbsPath: "/root/child2",
							Id:      "child2-id",
						},
					},
					Children: []*ds.TreeNode{
						{
							Info: &file.FileNode{
								GeneralNode: file.GeneralNode{
									AbsPath: "/root/child2/grandchild1",
									Tags:    []string{"grandchild-tag"},
									Id:      "grandchild1-id",
								},
							},
							Children: nil,
						},
						{
							Info: &file.FileNode{
								GeneralNode: file.GeneralNode{
									AbsPath: "/root/child2/grandchild2",
									Id:      "grandchild2-id",
								},
							},
							Children: nil,
						},
					},
				},
			},
		}

		err = rw.Write(original)
		require.NoError(t, err)

		read, err := rw.Read()
		require.NoError(t, err)
		require.NotNil(t, read)

		// Verify root
		originalRootFn := original.Info.(*file.FileNode)
		readRootFn := read.Info.(*file.FileNode)
		require.Equal(t, originalRootFn.AbsPath, readRootFn.AbsPath)
		require.Equal(t, originalRootFn.Tags, readRootFn.Tags)
		require.Equal(t, originalRootFn.Id, readRootFn.Id)

		// Verify children count
		require.Len(t, read.Children, 2)

		// Verify first child
		originalChild1Fn := original.Children[0].Info.(*file.FileNode)
		readChild1Fn := read.Children[0].Info.(*file.FileNode)
		require.Equal(t, originalChild1Fn.AbsPath, readChild1Fn.AbsPath)
		require.Equal(t, originalChild1Fn.Tags, readChild1Fn.Tags)
		require.Equal(t, originalChild1Fn.Id, readChild1Fn.Id)
		require.Empty(t, read.Children[0].Children)

		// Verify second child
		originalChild2Fn := original.Children[1].Info.(*file.FileNode)
		readChild2Fn := read.Children[1].Info.(*file.FileNode)
		require.Equal(t, originalChild2Fn.AbsPath, readChild2Fn.AbsPath)
		require.Equal(t, originalChild2Fn.Id, readChild2Fn.Id)

		// Verify grandchildren
		require.Len(t, read.Children[1].Children, 2)
		originalGrandchild1Fn := original.Children[1].Children[0].Info.(*file.FileNode)
		readGrandchild1Fn := read.Children[1].Children[0].Info.(*file.FileNode)
		require.Equal(t, originalGrandchild1Fn.AbsPath, readGrandchild1Fn.AbsPath)
		require.Equal(t, originalGrandchild1Fn.Tags, readGrandchild1Fn.Tags)
		require.Equal(t, originalGrandchild1Fn.Id, readGrandchild1Fn.Id)

		originalGrandchild2Fn := original.Children[1].Children[1].Info.(*file.FileNode)
		readGrandchild2Fn := read.Children[1].Children[1].Info.(*file.FileNode)
		require.Equal(t, originalGrandchild2Fn.AbsPath, readGrandchild2Fn.AbsPath)
		require.Equal(t, originalGrandchild2Fn.Id, readGrandchild2Fn.Id)
	})

	t.Run("write and read Google Drive node", func(t *testing.T) {
		dir := t.TempDir()
		dataFilePath := filepath.Join(dir, "data.json")

		rw, err := NewFileStorageRW(dataFilePath)
		require.NoError(t, err)

		original := &ds.TreeNode{
			Info: &file.FileNode{
				GeneralNode: file.GeneralNode{
					AbsPath: file.GDrivePathPrefix + "MyFolder/document.txt",
					Tags:    []string{"gdrive", "document"},
					Id:      "drive-node-id",
				},
				DriveId: "abc123def456",
			},
			Children: nil,
		}

		err = rw.Write(original)
		require.NoError(t, err)

		read, err := rw.Read()
		require.NoError(t, err)
		require.NotNil(t, read)

		originalFn := original.Info.(*file.FileNode)
		readFn := read.Info.(*file.FileNode)

		require.Equal(t, originalFn.AbsPath, readFn.AbsPath)
		require.Equal(t, originalFn.Tags, readFn.Tags)
		require.Equal(t, originalFn.Id, readFn.Id)
		require.Equal(t, originalFn.DriveId, readFn.DriveId)
	})

	t.Run("write and read node with empty tags", func(t *testing.T) {
		dir := t.TempDir()
		dataFilePath := filepath.Join(dir, "data.json")

		rw, err := NewFileStorageRW(dataFilePath)
		require.NoError(t, err)

		original := &ds.TreeNode{
			Info: &file.FileNode{
				GeneralNode: file.GeneralNode{
					AbsPath: "/empty/tags",
					Tags:    []string{},
					Id:      "",
				},
			},
			Children: nil,
		}

		err = rw.Write(original)
		require.NoError(t, err)

		read, err := rw.Read()
		require.NoError(t, err)
		require.NotNil(t, read)

		originalFn := original.Info.(*file.FileNode)
		readFn := read.Info.(*file.FileNode)

		require.Equal(t, originalFn.AbsPath, readFn.AbsPath)
		require.Equal(t, originalFn.Tags, readFn.Tags)
		require.Equal(t, originalFn.Id, readFn.Id)
	})
}
