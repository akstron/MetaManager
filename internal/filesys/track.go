package filesys

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/heroku/self/MetaManager/internal/cmderror"
	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/file"
	contextrepo "github.com/heroku/self/MetaManager/internal/repository/filesys"
	"github.com/heroku/self/MetaManager/internal/services"
)

type Tracker interface {
	// Tracks a path and returns a tree node
	Track(path string) (*ds.TreeNode, error)
}

// GDriveTracker tracks Google Drive paths using a GDriveServiceInterface.
// This allows for dependency injection and mocking in tests.
type GDriveTracker struct {
	svc services.GDriveServiceInterface
}

// NewGDriveTracker creates a new GDriveTracker with the given GDriveServiceInterface.
func NewGDriveTracker(svc services.GDriveServiceInterface) *GDriveTracker {
	return &GDriveTracker{svc: svc}
}

func (g *GDriveTracker) Track(path string) (*ds.TreeNode, error) {
	if len(path) == 0 {
		return nil, &cmderror.InvalidPath{}
	}

	scanner := NewGDriveScanner(g.svc)
	ctx := context.Background()

	if path[len(path)-1] != '*' {
		// Recursive tracking: remove '*' and scan
		// rootDirPath := path[0 : len(path)-1]
		// Normalize the path (remove gdrive: prefix if present, handle relative paths)
		// drivePath, _ := scanner.NormalizeTrackPath(rootDirPath)
		// return scanner.TrackGDrive(ctx, drivePath, true)
		rootDirPath := path
		// if isPathPrefixOrEqual(rootDirPath, file.GDrivePathPrefix) {
		rootDirPath = strings.Trim(rootDirPath, "/")
		rootDirPath = file.GDrivePathPrefix + rootDirPath
		// }
		return file.CreateTreeNodeFromPath(rootDirPath)
	}

	// Non-recursive tracking: just create a node for the path
	drivePath, _ := scanner.NormalizeTrackPath(path)
	return scanner.TrackGDrive(ctx, drivePath, true)
}

// LocalTracker tracks local filesystem paths.
type LocalTracker struct {
	scanner *UnixFileSystemScanner
}

// NewLocalTracker creates a new LocalTracker.
func NewLocalTracker() *LocalTracker {
	return &LocalTracker{
		scanner: NewUnixFileSystemScanner(),
	}
}

func (l *LocalTracker) Track(path string) (*ds.TreeNode, error) {
	if len(path) == 0 {
		return nil, &cmderror.InvalidPath{}
	}

	if path[len(path)-1] == '*' {
		// Recursive tracking: remove '*' and scan
		rootDirPath := path[0 : len(path)-1]
		absPath, err := filepath.Abs(rootDirPath)
		if err != nil {
			return nil, err
		}
		return l.scanner.Scan(absPath)
	}

	// Non-recursive tracking: just create a node for the path
	return file.CreateTreeNodeFromPath(path)
}

// ContextAwareTracker creates the appropriate tracker based on the current context type.
type ContextAwareTracker struct {
	cxtRepo contextrepo.ContextRepository
}

func NewContextAwareTracker(cxtRepo contextrepo.ContextRepository) *ContextAwareTracker {
	return &ContextAwareTracker{cxtRepo: cxtRepo}
}

func (c *ContextAwareTracker) Track(path string) (*ds.TreeNode, error) {
	if len(path) == 0 {
		return nil, &cmderror.InvalidPath{}
	}

	ctxName, err := c.cxtRepo.GetContext()
	if err != nil {
		return nil, err
	}
	contextType, err := c.cxtRepo.GetContextType(ctxName)
	if err != nil {
		return nil, err
	}

	var tracker Tracker
	switch contextType {
	case contextrepo.TypeGDrive:
		// For GDrive, we need to get the service
		// In production, this calls services.GetGDriveService()
		// In tests, this can be mocked by creating a GDriveTracker directly
		svc, err := services.GetGDriveService(context.Background())
		if err != nil {
			return nil, err
		}
		tracker = NewGDriveTracker(svc)
	case contextrepo.TypeLocal:
		tracker = NewLocalTracker()
	default:
		return nil, &cmderror.InvalidOperation{}
	}

	return tracker.Track(path)
}
