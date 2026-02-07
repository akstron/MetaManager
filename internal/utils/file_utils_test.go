package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindMMDirPath(t *testing.T) {
	dir := t.TempDir()
	os.Setenv("MM_TEST_CONTEXT_DIR", dir)
	defer os.Unsetenv("MM_TEST_CONTEXT_DIR")

	parentDir, err := GetAppDataDir()
	require.NoError(t, err)
	require.Equal(t, filepath.Join(dir, MMDirName), parentDir)

	appDir, err := GetAppDataDirForContext("myctx")
	require.NoError(t, err)
	require.Equal(t, filepath.Join(dir, MMDirName, "myctx"), appDir)

	// Before creating .mm/myctx, FindMMDirPath returns false
	found, path, err := FindMMDirPath("myctx")
	require.NoError(t, err)
	require.False(t, found)
	require.Empty(t, path)

	require.NoError(t, os.MkdirAll(appDir, 0755))
	found, path, err = FindMMDirPath("myctx")
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, appDir, path)
}
