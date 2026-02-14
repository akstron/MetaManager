package filesys

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/heroku/self/MetaManager/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestContextRepositoryImpl_gdriveCwdPath(t *testing.T) {
	dir := t.TempDir()
	os.Setenv("MM_TEST_CONTEXT_DIR", dir)
	defer os.Unsetenv("MM_TEST_CONTEXT_DIR")

	store := NewContextRepositoryImpl(nil)

	t.Run("empty context name returns error", func(t *testing.T) {
		_, err := store.gdriveCwdPath("")
		require.Error(t, err)
		require.Contains(t, err.Error(), "context name cannot be empty")
	})

	t.Run("returns path under .mm/<contextName>/gdrive_cwd", func(t *testing.T) {
		path, err := store.gdriveCwdPath("myctx")
		require.NoError(t, err)
		expected := filepath.Join(dir, utils.MMDirName, "myctx", gdriveCwdFileName)
		require.Equal(t, expected, path)
	})
}
