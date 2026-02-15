package file

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileNodeJSONSerializer_InfoUnmarshal(t *testing.T) {
	serializer := FileNodeJSONSerializer{}

	t.Run("local file node with all fields", func(t *testing.T) {
		info := map[string]interface{}{
			"AbsPath": "/path/to/file",
			"Tags":    []interface{}{"tag1", "tag2"},
			"Id":      "node-id-123",
			"DriveId": "",
		}

		result, err := serializer.InfoUnmarshal(info)
		require.NoError(t, err)
		require.NotNil(t, result)

		fn, ok := result.(*FileNode)
		require.True(t, ok, "result should be *FileNode")
		require.Equal(t, "/path/to/file", fn.AbsPath)
		require.Equal(t, []string{"tag1", "tag2"}, fn.Tags)
		require.Equal(t, "node-id-123", fn.Id)
		require.Equal(t, "", fn.DriveId)
		require.Equal(t, "FILE", fn.Name())
	})

	t.Run("Google Drive node with DriveId", func(t *testing.T) {
		info := map[string]interface{}{
			"AbsPath": "gdrive:/Folder/file.txt",
			"Tags":    []interface{}{"gdrive", "document"},
			"Id":      "drive-node-456",
			"DriveId": "1a2b3c4d5e6f7g8h",
		}

		result, err := serializer.InfoUnmarshal(info)
		require.NoError(t, err)
		require.NotNil(t, result)

		fn, ok := result.(*FileNode)
		require.True(t, ok, "result should be *FileNode")
		require.Equal(t, "gdrive:/Folder/file.txt", fn.AbsPath)
		require.Equal(t, []string{"gdrive", "document"}, fn.Tags)
		require.Equal(t, "drive-node-456", fn.Id)
		require.Equal(t, "1a2b3c4d5e6f7g8h", fn.DriveId)
	})

	t.Run("minimal node with only AbsPath", func(t *testing.T) {
		info := map[string]interface{}{
			"AbsPath": "/root",
		}

		result, err := serializer.InfoUnmarshal(info)
		require.NoError(t, err)
		require.NotNil(t, result)

		fn, ok := result.(*FileNode)
		require.True(t, ok, "result should be *FileNode")
		require.Equal(t, "/root", fn.AbsPath)
		require.Nil(t, fn.Tags)
		require.Equal(t, "", fn.Id)
		require.Equal(t, "", fn.DriveId)
	})

	t.Run("node with empty tags slice", func(t *testing.T) {
		info := map[string]interface{}{
			"AbsPath": "/path/to/dir",
			"Tags":    []interface{}{},
			"Id":      "empty-tags-id",
		}

		result, err := serializer.InfoUnmarshal(info)
		require.NoError(t, err)
		require.NotNil(t, result)

		fn, ok := result.(*FileNode)
		require.True(t, ok, "result should be *FileNode")
		require.Equal(t, "/path/to/dir", fn.AbsPath)
		require.Equal(t, []string{}, fn.Tags)
		require.Equal(t, "empty-tags-id", fn.Id)
	})

	t.Run("node with single tag", func(t *testing.T) {
		info := map[string]interface{}{
			"AbsPath": "/single/tag/path",
			"Tags":    []interface{}{"important"},
		}

		result, err := serializer.InfoUnmarshal(info)
		require.NoError(t, err)
		require.NotNil(t, result)

		fn, ok := result.(*FileNode)
		require.True(t, ok, "result should be *FileNode")
		require.Equal(t, "/single/tag/path", fn.AbsPath)
		require.Equal(t, []string{"important"}, fn.Tags)
	})

	t.Run("root path node", func(t *testing.T) {
		info := map[string]interface{}{
			"AbsPath": "/",
			"Id":      "root-id",
		}

		result, err := serializer.InfoUnmarshal(info)
		require.NoError(t, err)
		require.NotNil(t, result)

		fn, ok := result.(*FileNode)
		require.True(t, ok, "result should be *FileNode")
		require.Equal(t, "/", fn.AbsPath)
		require.Equal(t, "root-id", fn.Id)
	})

	t.Run("gdrive root path", func(t *testing.T) {
		info := map[string]interface{}{
			"AbsPath": GDrivePathPrefix,
			"Id":      "gdrive-root",
		}

		result, err := serializer.InfoUnmarshal(info)
		require.NoError(t, err)
		require.NotNil(t, result)

		fn, ok := result.(*FileNode)
		require.True(t, ok, "result should be *FileNode")
		require.Equal(t, GDrivePathPrefix, fn.AbsPath)
		require.Equal(t, "gdrive-root", fn.Id)
	})

	t.Run("handles missing optional fields", func(t *testing.T) {
		info := map[string]interface{}{
			"AbsPath": "/minimal/path",
		}

		result, err := serializer.InfoUnmarshal(info)
		require.NoError(t, err)
		require.NotNil(t, result)

		fn, ok := result.(*FileNode)
		require.True(t, ok, "result should be *FileNode")
		require.Equal(t, "/minimal/path", fn.AbsPath)
		// Tags, Id, and DriveId should be zero values
		require.Nil(t, fn.Tags)
		require.Equal(t, "", fn.Id)
		require.Equal(t, "", fn.DriveId)
	})
}
