package filesys

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/heroku/self/MetaManager/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRepo(t *testing.T) *ContextRepositoryImpl {
	dir := t.TempDir()
	return NewContextRepositoryImpl(func() (string, error) { return dir, nil })
}

func TestNewContextRepositoryImpl(t *testing.T) {
	// With custom baseDir
	dir := t.TempDir()
	repo := NewContextRepositoryImpl(func() (string, error) { return dir, nil })
	require.NotNil(t, repo)
	path, err := repo.ContextsJSONPath()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(dir, contextsJSONFileName), path)

	// With nil baseDir uses default (MM_TEST_CONTEXT_DIR or executable dir)
	os.Setenv("MM_TEST_CONTEXT_DIR", dir)
	defer os.Unsetenv("MM_TEST_CONTEXT_DIR")
	repo2 := NewContextRepositoryImpl(nil)
	require.NotNil(t, repo2)
	path2, err := repo2.ContextsJSONPath()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(dir, contextsJSONFileName), path2)
}

func TestContextRepositoryImpl_ContextFilePath_ContextsJSONPath(t *testing.T) {
	dir := t.TempDir()
	repo := NewContextRepositoryImpl(func() (string, error) { return dir, nil })

	ctxPath, err := repo.ContextFilePath()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(dir, contextFileName), ctxPath)

	jsonPath, err := repo.ContextsJSONPath()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(dir, contextsJSONFileName), jsonPath)
}

func TestContextRepositoryImpl_LoadContexts(t *testing.T) {
	repo := testRepo(t)

	// No file -> nil
	entries, err := repo.LoadContexts()
	require.NoError(t, err)
	assert.Nil(t, entries)

	// Write contexts.json
	path, err := repo.ContextsJSONPath()
	require.NoError(t, err)
	require.NoError(t, utils.WriteJSON(path, ContextsFile{
		Contexts: []ContextEntry{
			{Name: "a", Type: TypeLocal},
			{Name: "b", Type: TypeGDrive},
		},
	}, true))

	entries, err = repo.LoadContexts()
	require.NoError(t, err)
	require.Len(t, entries, 2)
	assert.Equal(t, "a", entries[0].Name)
	assert.Equal(t, TypeLocal, entries[0].Type)
	assert.Equal(t, "b", entries[1].Name)
	assert.Equal(t, TypeGDrive, entries[1].Type)
}

func TestContextRepositoryImpl_Create(t *testing.T) {
	repo := testRepo(t)

	err := repo.Create("work", TypeLocal)
	require.NoError(t, err)

	entries, err := repo.LoadContexts()
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, "work", entries[0].Name)
	assert.Equal(t, TypeLocal, entries[0].Type)

	// Second context, sorted by name
	err = repo.Create("gdrive", TypeGDrive)
	require.NoError(t, err)
	entries, err = repo.LoadContexts()
	require.NoError(t, err)
	require.Len(t, entries, 2)
	assert.Equal(t, "gdrive", entries[0].Name)
	assert.Equal(t, "work", entries[1].Name)

	// Duplicate name
	err = repo.Create("work", TypeLocal)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestContextRepositoryImpl_Create_Validation(t *testing.T) {
	repo := testRepo(t)

	err := repo.Create("", TypeLocal)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")

	err = repo.Create("x", "invalid")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "type must be")
}

func TestContextRepositoryImpl_AddContext(t *testing.T) {
	repo := testRepo(t)

	err := repo.AddContext("alpha", TypeLocal)
	require.NoError(t, err)
	err = repo.AddContext("beta", TypeGDrive)
	require.NoError(t, err)

	entries, err := repo.LoadContexts()
	require.NoError(t, err)
	require.Len(t, entries, 2)
	assert.Equal(t, "alpha", entries[0].Name)
	assert.Equal(t, "beta", entries[1].Name)
}

func TestContextRepositoryImpl_SetCurrent(t *testing.T) {
	repo := testRepo(t)
	require.NoError(t, repo.Create("work", TypeLocal))
	require.NoError(t, repo.Create("gdrive", TypeGDrive))

	err := repo.SetCurrent("work")
	require.NoError(t, err)
	ctxPath, err := repo.ContextFilePath()
	require.NoError(t, err)
	b, err := os.ReadFile(ctxPath)
	require.NoError(t, err)
	assert.Equal(t, "work", string(b))

	err = repo.SetCurrent("nonexistent")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestContextRepositoryImpl_SetCurrent_Validation(t *testing.T) {
	repo := testRepo(t)

	err := repo.SetCurrent("")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestContextRepositoryImpl_SetCurrent_NormalizesName(t *testing.T) {
	repo := testRepo(t)
	require.NoError(t, repo.Create("work", TypeLocal))

	err := repo.SetCurrent("  WORK  ")
	require.NoError(t, err)
	name, err := repo.GetContext()
	require.NoError(t, err)
	assert.Equal(t, "work", name)
}

func TestContextRepositoryImpl_GetContext(t *testing.T) {
	repo := testRepo(t)

	// No file, no env
	name, err := repo.GetContext()
	require.NoError(t, err)
	assert.Empty(t, name)

	// From file
	require.NoError(t, repo.Create("local1", TypeLocal))
	require.NoError(t, repo.SetCurrent("local1"))
	name, err = repo.GetContext()
	require.NoError(t, err)
	assert.Equal(t, "local1", name)

	// Env overrides file
	os.Setenv(ContextEnvVar, "  ENVCTX  ")
	defer os.Unsetenv(ContextEnvVar)
	name, err = repo.GetContext()
	require.NoError(t, err)
	assert.Equal(t, "envctx", name)
}

func TestContextRepositoryImpl_GetContextType(t *testing.T) {
	repo := testRepo(t)
	path, err := repo.ContextsJSONPath()
	require.NoError(t, err)
	require.NoError(t, utils.WriteJSON(path, ContextsFile{
		Contexts: []ContextEntry{
			{Name: "local-ctx", Type: TypeLocal},
			{Name: "drive-ctx", Type: TypeGDrive},
		},
	}, true))

	typ, err := repo.GetContextType("local-ctx")
	require.NoError(t, err)
	assert.Equal(t, TypeLocal, typ)

	typ, err = repo.GetContextType("drive-ctx")
	require.NoError(t, err)
	assert.Equal(t, TypeGDrive, typ)

	typ, err = repo.GetContextType("missing")
	require.NoError(t, err)
	assert.Empty(t, typ)
}

func TestContextRepositoryImpl_AsContextRepository(t *testing.T) {
	repo := testRepo(t)
	var _ ContextRepository = repo

	// Use via interface
	var cr ContextRepository = repo
	require.NoError(t, cr.Create("via-interface", TypeLocal))
	entries, err := cr.LoadContexts()
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, "via-interface", entries[0].Name)
}

func TestContextRepositoryImpl_BaseDirError(t *testing.T) {
	errDir := func() (string, error) { return "", os.ErrNotExist }
	repo := NewContextRepositoryImpl(errDir)

	_, err := repo.ContextFilePath()
	require.Error(t, err)
	assert.ErrorIs(t, err, os.ErrNotExist)

	_, err = repo.ContextsJSONPath()
	require.Error(t, err)

	_, err = repo.LoadContexts()
	require.Error(t, err)
}

func TestContextRepositoryImpl_GetGDriveCwd_SetGDriveCwd(t *testing.T) {
	repo := testRepo(t)

	// Unset -> empty (caller treats as "/")
	cwd, err := repo.GetGDriveCwd()
	require.NoError(t, err)
	assert.Empty(t, cwd)

	require.NoError(t, repo.SetGDriveCwd("/Foo/Bar"))
	cwd, err = repo.GetGDriveCwd()
	require.NoError(t, err)
	assert.Equal(t, "/Foo/Bar", cwd)

	require.NoError(t, repo.SetGDriveCwd("/"))
	cwd, err = repo.GetGDriveCwd()
	require.NoError(t, err)
	assert.Equal(t, "/", cwd)

	// Env overrides file
	os.Setenv(GDriveCwdEnvVar, "/Env/Path")
	defer os.Unsetenv(GDriveCwdEnvVar)
	cwd, err = repo.GetGDriveCwd()
	require.NoError(t, err)
	assert.Equal(t, "/Env/Path", cwd)
}

func TestResolveGDrivePath(t *testing.T) {
	tests := []struct {
		cwd    string
		target string
		want   string
	}{
		{"", "", "/"},
		{"/", "", "/"},
		{"/", ".", "/"},
		{"/Foo", ".", "/Foo"},
		{"/", "/Bar", "/Bar"},
		{"/Foo", "/Bar", "/Bar"},
		{"/Foo", "Sub", "/Foo/Sub"},
		{"/Foo/Bar", "Sub", "/Foo/Bar/Sub"},
		{"/Foo/Bar", "..", "/Foo"},
		{"/Foo/Bar", "../Baz", "/Foo/Baz"},
	}
	for _, tt := range tests {
		got := ResolveGDrivePath(tt.cwd, tt.target)
		assert.Equal(t, tt.want, got, "cwd=%q target=%q", tt.cwd, tt.target)
	}
}
