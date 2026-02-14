package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/heroku/self/MetaManager/internal/googleauth"
	"github.com/heroku/self/MetaManager/internal/utils"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const (
	// DriveRootID is the fixed ID for the root of "My Drive" in Google Drive.
	DriveRootID = "root"
	// DriveFolderMimeType is the MIME type for Drive folders.
	DriveFolderMimeType = "application/vnd.google-apps.folder"
	// GoogleTokenFileName is the name of the file where the OAuth token is stored (in base dir).
	GoogleTokenFileName = "google_token.json"
)

// embeddedCredentials is set by main when credentials.json is embedded in the binary (via SetEmbeddedCredentials).
var embeddedCredentials []byte

// SetEmbeddedCredentials sets the credentials JSON embedded in the binary. Call from main after embedding credentials.json.
func SetEmbeddedCredentials(b []byte) {
	embeddedCredentials = b
}

// EmbeddedCredentials returns the embedded credentials JSON (for use by login flow). Returns nil if not set.
func EmbeddedCredentials() []byte {
	return embeddedCredentials
}

// GDriveService authenticates with Google Drive using a stored token and lists directory structure.
type GDriveService struct {
	svc *drive.Service
}

// RootEntry represents a single file or folder in Google Drive (used for any directory listing).
type RootEntry struct {
	Id       string // Drive file ID
	Name     string
	IsFolder bool
	MimeType string
}

// NewGDriveService creates a Drive API client using the given OAuth config and token.
// The token is used to build a TokenSource that can refresh when expired.
func NewGDriveService(ctx context.Context, config *oauth2.Config, token *oauth2.Token) (*GDriveService, error) {
	ts := config.TokenSource(ctx, token)
	svc, err := drive.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		return nil, fmt.Errorf("create drive service: %w", err)
	}
	return &GDriveService{svc: svc}, nil
}

// NewGDriveServiceFromTokenPath creates a Drive client by loading the token from path
// and using the provided credentials (e.g. from embedded credentials.json) for refresh.
func NewGDriveServiceFromTokenPath(ctx context.Context, tokenPath string, credentialsJSON []byte) (*GDriveService, error) {
	config, err := googleauth.LoadConfigFromBytes(credentialsJSON)
	if err != nil {
		return nil, fmt.Errorf("load credentials: %w", err)
	}
	token, err := googleauth.LoadToken(tokenPath)
	if err != nil {
		return nil, fmt.Errorf("load token from %q: %w", tokenPath, err)
	}
	return NewGDriveService(ctx, config, token)
}

// TokenPath returns the path to the Google OAuth token file (in base dir: executable dir or MM_TEST_CONTEXT_DIR).
func TokenPath() (string, error) {
	baseDir, err := utils.GetBaseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(baseDir, GoogleTokenFileName), nil
}

// GetGDriveService discovers the token path and uses embedded credentials to return a GDrive service.
// Call SetEmbeddedCredentials from main first. The token file must exist (run login first).
func GetGDriveService(ctx context.Context) (*GDriveService, error) {
	if len(embeddedCredentials) == 0 {
		return nil, fmt.Errorf("no credentials; rebuild the binary with credentials.json for Drive")
	}
	tokenPath, err := TokenPath()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(tokenPath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("token not found at %q; run login first for Drive", tokenPath)
		}
		return nil, err
	}
	return NewGDriveServiceFromTokenPath(ctx, tokenPath, embeddedCredentials)
}

// ListFolder returns the immediate children of the given folder (by Drive file ID), excluding trashed items.
// Use DriveRootID for the root of "My Drive". Entries are sorted by name (folders first, then files).
func (g *GDriveService) ListFolder(ctx context.Context, folderID string) ([]RootEntry, error) {
	parentID := folderID
	if parentID == "" {
		parentID = DriveRootID
	}
	q := fmt.Sprintf("%q in parents and trashed = false", parentID)
	call := g.svc.Files.List().
		Q(q).
		Fields("nextPageToken, files(id, name, mimeType)").
		PageSize(1000)
	var all []RootEntry
	for {
		r, err := call.Context(ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("drive files.list: %w", err)
		}
		for _, f := range r.Files {
			all = append(all, RootEntry{
				Id:       f.Id,
				Name:     f.Name,
				IsFolder: f.MimeType == DriveFolderMimeType,
				MimeType: f.MimeType,
			})
		}
		if r.NextPageToken == "" {
			break
		}
		call = call.PageToken(r.NextPageToken)
	}
	sort.Slice(all, func(i, j int) bool {
		if all[i].IsFolder != all[j].IsFolder {
			return all[i].IsFolder
		}
		return all[i].Name < all[j].Name
	})
	return all, nil
}

// ListRoot returns the immediate children of the Drive root (My Drive). Same as ListFolder(ctx, DriveRootID).
func (g *GDriveService) ListRoot(ctx context.Context) ([]RootEntry, error) {
	return g.ListFolder(ctx, DriveRootID)
}

// ResolvePath resolves a path like "/" or "/Folder1" or "/Folder1/SubFolder" to the Drive folder ID at that path.
// Path is slash-separated; leading and trailing slashes are ignored. "/" or "" returns DriveRootID.
// Returns error if any segment is not found or is not a folder.
func (g *GDriveService) ResolvePath(ctx context.Context, path string) (folderID string, err error) {
	path = strings.Trim(path, "/")
	if path == "" {
		return DriveRootID, nil
	}
	parts := strings.Split(path, "/")
	parentID := DriveRootID
	for i, name := range parts {
		entries, err := g.ListFolder(ctx, parentID)
		if err != nil {
			return "", fmt.Errorf("list %q: %w", path, err)
		}
		var found *RootEntry
		for j := range entries {
			if entries[j].Name == name {
				found = &entries[j]
				break
			}
		}
		if found == nil {
			return "", fmt.Errorf("path not found: %q (no %q in %q)", path, name, strings.Join(parts[:i], "/"))
		}
		if !found.IsFolder {
			return "", fmt.Errorf("not a folder: %q", name)
		}
		parentID = found.Id
	}
	return parentID, nil
}

// ListAtPath lists the contents of the folder at the given path.
// Path is like "/" for root, "/Folder1", or "/Folder1/SubFolder". Same semantics as ResolvePath.
func (g *GDriveService) ListAtPath(ctx context.Context, path string) ([]RootEntry, error) {
	folderID, err := g.ResolvePath(ctx, path)
	if err != nil {
		return nil, err
	}
	return g.ListFolder(ctx, folderID)
}
