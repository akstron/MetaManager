package services

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetEmbeddedCredentials_EmbeddedCredentials(t *testing.T) {
	// Restore original so we don't affect other tests
	orig := embeddedCredentials
	defer func() { embeddedCredentials = orig }()

	require.Nil(t, EmbeddedCredentials())
	SetEmbeddedCredentials([]byte("{}"))
	require.Equal(t, []byte("{}"), EmbeddedCredentials())
	SetEmbeddedCredentials([]byte(`{"installed":{"client_id":"x"}}`))
	require.Equal(t, []byte(`{"installed":{"client_id":"x"}}`), EmbeddedCredentials())
	SetEmbeddedCredentials(nil)
	require.Nil(t, EmbeddedCredentials())
}

func TestTokenPath(t *testing.T) {
	dir := t.TempDir()
	os.Setenv("MM_TEST_CONTEXT_DIR", dir)
	defer os.Unsetenv("MM_TEST_CONTEXT_DIR")

	path, err := TokenPath()
	require.NoError(t, err)
	require.Equal(t, filepath.Join(dir, GoogleTokenFileName), path)
}

func TestGetGDriveService_NoCredentials(t *testing.T) {
	orig := embeddedCredentials
	defer func() { embeddedCredentials = orig }()
	embeddedCredentials = nil

	ctx := context.Background()
	_, err := GetGDriveService(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no credentials")
}

func TestGetGDriveService_TokenFileMissing(t *testing.T) {
	dir := t.TempDir()
	os.Setenv("MM_TEST_CONTEXT_DIR", dir)
	defer os.Unsetenv("MM_TEST_CONTEXT_DIR")

	orig := embeddedCredentials
	defer func() { embeddedCredentials = orig }()
	embeddedCredentials = []byte(`{"installed":{"client_id":"x","client_secret":"y","auth_uri":"https://accounts.google.com/o/oauth2/auth","token_uri":"https://oauth2.googleapis.com/token"}}`)

	ctx := context.Background()
	_, err := GetGDriveService(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "token not found")
	require.Contains(t, err.Error(), "run login first")
}

func TestNewGDriveServiceFromTokenPath_EmptyCredentials(t *testing.T) {
	ctx := context.Background()
	_, err := NewGDriveServiceFromTokenPath(ctx, "/nonexistent/token.json", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "load credentials")
}

func TestNewGDriveServiceFromTokenPath_InvalidCredentialsJSON(t *testing.T) {
	ctx := context.Background()
	_, err := NewGDriveServiceFromTokenPath(ctx, "/nonexistent/token.json", []byte("not json"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "load credentials")
}

func TestNewGDriveServiceFromTokenPath_TokenFileNotFound(t *testing.T) {
	// Valid OAuth client config shape so LoadConfigFromBytes passes and we hit LoadToken (missing file).
	creds := []byte(`{"installed":{"client_id":"x","client_secret":"y","auth_uri":"https://accounts.google.com/o/oauth2/auth","token_uri":"https://oauth2.googleapis.com/token","redirect_uris":["urn:ietf:wg:oauth:2.0:oob"]}}`)
	ctx := context.Background()
	_, err := NewGDriveServiceFromTokenPath(ctx, "/nonexistent/path/google_token.json", creds)
	require.Error(t, err)
	require.Contains(t, err.Error(), "load token")
}

func TestConstants(t *testing.T) {
	require.Equal(t, "root", DriveRootID)
	require.Equal(t, "application/vnd.google-apps.folder", DriveFolderMimeType)
	require.Equal(t, "google_token.json", GoogleTokenFileName)
}
