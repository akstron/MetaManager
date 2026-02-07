package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/heroku/self/MetaManager/internal/repository/filesys"
	"github.com/heroku/self/MetaManager/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testStore(t *testing.T) *filesys.ContextRepositoryImpl {
	dir := t.TempDir()
	return filesys.NewContextRepositoryImpl(func() (string, error) { return dir, nil })
}

func TestContextCreateAndSet(t *testing.T) {
	store := testStore(t)

	// Create context "work" (local)
	err := store.Create("work", filesys.TypeLocal)
	require.NoError(t, err)

	// contexts.json should exist with one entry
	path, err := store.ContextsJSONPath()
	require.NoError(t, err)
	var f filesys.ContextsFile
	err = utils.ReadJSON(path, &f)
	require.NoError(t, err)
	require.Len(t, f.Contexts, 1)
	assert.Equal(t, "work", f.Contexts[0].Name)
	assert.Equal(t, filesys.TypeLocal, f.Contexts[0].Type)

	// Create another context
	err = store.Create("gdrive", filesys.TypeGDrive)
	require.NoError(t, err)
	err = utils.ReadJSON(path, &f)
	require.NoError(t, err)
	require.Len(t, f.Contexts, 2)
	assert.Equal(t, "gdrive", f.Contexts[0].Name)
	assert.Equal(t, "work", f.Contexts[1].Name)

	// Duplicate name should error
	err = store.Create("work", filesys.TypeLocal)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")

	// Set current context to work
	err = store.SetCurrent("work")
	require.NoError(t, err)
	contextPath, err := store.ContextFilePath()
	require.NoError(t, err)
	b, err := os.ReadFile(contextPath)
	require.NoError(t, err)
	assert.Equal(t, "work", string(b))

	// Set to non-existent context should error
	err = store.SetCurrent("nonexistent")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestContextCreateValidation(t *testing.T) {
	store := testStore(t)

	// Empty name
	err := store.Create("", filesys.TypeLocal)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")

	// Invalid type
	err = store.Create("x", "invalid")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "type must be")
}

func TestContextSetValidation(t *testing.T) {
	store := testStore(t)

	// Empty name
	err := store.SetCurrent("")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestGetContext(t *testing.T) {
	store := testStore(t)

	// No file and no env -> empty
	name, err := store.GetContext()
	require.NoError(t, err)
	assert.Empty(t, name)

	// With env set (env is process-wide, so we test store's use of it)
	os.Setenv(filesys.ContextEnvVar, "  WORK  ")
	defer os.Unsetenv(filesys.ContextEnvVar)
	name, err = store.GetContext()
	require.NoError(t, err)
	assert.Equal(t, "work", name)

	os.Unsetenv(filesys.ContextEnvVar)

	// With context file
	require.NoError(t, store.Create("gdrive", filesys.TypeGDrive))
	require.NoError(t, store.SetCurrent("gdrive"))
	name, err = store.GetContext()
	require.NoError(t, err)
	assert.Equal(t, "gdrive", name)

	// Env overrides file
	os.Setenv(filesys.ContextEnvVar, "local")
	defer os.Unsetenv(filesys.ContextEnvVar)
	name, err = store.GetContext()
	require.NoError(t, err)
	assert.Equal(t, "local", name)
}

func TestGetContexts(t *testing.T) {
	store := testStore(t)

	// No file -> nil
	entries, err := store.LoadContexts()
	require.NoError(t, err)
	assert.Nil(t, entries)

	// Write contexts.json
	path, err := store.ContextsJSONPath()
	require.NoError(t, err)
	require.NoError(t, utils.WriteJSON(path, filesys.ContextsFile{
		Contexts: []filesys.ContextEntry{
			{Name: "a", Type: filesys.TypeLocal},
			{Name: "b", Type: filesys.TypeGDrive},
		},
	}, true))

	entries, err = store.LoadContexts()
	require.NoError(t, err)
	require.Len(t, entries, 2)
	assert.Equal(t, "a", entries[0].Name)
	assert.Equal(t, filesys.TypeLocal, entries[0].Type)
	assert.Equal(t, "b", entries[1].Name)
	assert.Equal(t, filesys.TypeGDrive, entries[1].Type)
}

func TestGetContextType(t *testing.T) {
	store := testStore(t)

	path, err := store.ContextsJSONPath()
	require.NoError(t, err)
	require.NoError(t, utils.WriteJSON(path, filesys.ContextsFile{
		Contexts: []filesys.ContextEntry{
			{Name: "local-ctx", Type: filesys.TypeLocal},
			{Name: "drive-ctx", Type: filesys.TypeGDrive},
		},
	}, true))

	typ, err := store.GetContextType("local-ctx")
	require.NoError(t, err)
	assert.Equal(t, filesys.TypeLocal, typ)

	typ, err = store.GetContextType("drive-ctx")
	require.NoError(t, err)
	assert.Equal(t, filesys.TypeGDrive, typ)

	typ, err = store.GetContextType("missing")
	require.NoError(t, err)
	assert.Empty(t, typ)
}

// TestPackageLevelGetters ensures GetContext/GetContexts/GetContextType use defaultStore (e.g. with MM_TEST_CONTEXT_DIR).
func TestPackageLevelGetters(t *testing.T) {
	dir := t.TempDir()
	os.Setenv("MM_TEST_CONTEXT_DIR", dir)
	defer os.Unsetenv("MM_TEST_CONTEXT_DIR")

	os.Unsetenv(filesys.ContextEnvVar)
	name, err := GetContext()
	require.NoError(t, err)
	assert.Empty(t, name)

	entries, err := GetContexts()
	require.NoError(t, err)
	assert.Nil(t, entries)

	// Write via default store path
	path := filepath.Join(dir, "contexts.json")
	require.NoError(t, utils.WriteJSON(path, filesys.ContextsFile{
		Contexts: []filesys.ContextEntry{{Name: "pkg", Type: filesys.TypeLocal}},
	}, true))
	entries, err = GetContexts()
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, "pkg", entries[0].Name)

	typ, err := GetContextType("pkg")
	require.NoError(t, err)
	assert.Equal(t, filesys.TypeLocal, typ)
}
