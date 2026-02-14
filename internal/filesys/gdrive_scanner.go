package filesys

import (
	"context"
	"path"
	"strings"

	"github.com/heroku/self/MetaManager/internal/cmderror"
	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/file"
	"github.com/heroku/self/MetaManager/internal/services"
	"github.com/sirupsen/logrus"
)

const maxTrackDepth = 50

// Drive shortcut MIME type: do not recurse into these (they point to other folders and can create cycles).
const driveShortcutMimeType = "application/vnd.google-apps.shortcut"

// GDriveScanner scans Google Drive directories.
type GDriveScanner struct {
	svc *services.GDriveService
}

// NewGDriveScanner creates a new GDriveScanner with the given service.
func NewGDriveScanner(svc *services.GDriveService) *GDriveScanner {
	return &GDriveScanner{svc: svc}
}

// createGDriveScanner creates a new GDriveScanner.
func createGDriveScanner() (*GDriveScanner, error) {
	svc, err := services.GetGDriveService(context.Background())
	if err != nil {
		return nil, err
	}
	return NewGDriveScanner(svc), nil
}

// Scan scans a Google Drive path and returns a tree node.
func (g *GDriveScanner) Scan(path string) (*ds.TreeNode, error) {
	ctx := context.Background()
	drivePath, recursive := g.NormalizeTrackPath(path)
	return g.TrackGDrive(ctx, drivePath, recursive)
}

// Make sure that GDriveScanner implements Scanner
var _ Scanner = (*GDriveScanner)(nil)

// TrackGDrive builds a subtree from a Google Drive path. Path is like "/" for root or "/Folder/SubFolder".
// If recursive is true, recurses into all subfolders (up to maxTrackDepth).
// Shortcuts or shared folders that create cycles are skipped to avoid infinite recursion.
func (g *GDriveScanner) TrackGDrive(ctx context.Context, drivePath string, recursive bool) (*ds.TreeNode, error) {
	logrus.Debugf("[track-gdrive] TrackGDrive start path=%q recursive=%v", drivePath, recursive)
	if g.svc == nil {
		return nil, &cmderror.InvalidOperation{}
	}
	drivePath = strings.Trim(drivePath, "/")
	folderID, err := g.svc.ResolvePath(ctx, drivePath)
	if err != nil {
		logrus.Debugf("[track-gdrive] ResolvePath error: %v", err)
		return nil, err
	}
	logrus.Debugf("[track-gdrive] resolved folderID=%q", folderID)
	baseVirtual := file.GDrivePathPrefix
	if drivePath != "" {
		baseVirtual = file.GDrivePathPrefix + drivePath
	}
	visited := make(map[string]bool)
	visited[folderID] = true // avoid cycling back to root
	return g.trackGDriveFolder(ctx, folderID, baseVirtual, recursive, 0, visited)
}

// trackGDriveFolder recursively tracks a Google Drive folder.
func (g *GDriveScanner) trackGDriveFolder(ctx context.Context, folderID, virtualPath string, recursive bool, depth int, visited map[string]bool) (*ds.TreeNode, error) {
	logrus.Debugf("[track-gdrive] trackGDriveFolder depth=%d folderID=%q path=%q", depth, folderID, virtualPath)

	if depth > maxTrackDepth {
		logrus.Debugf("[track-gdrive] max depth %d exceeded, stopping", maxTrackDepth)
		return nil, &cmderror.InvalidOperation{}
	}
	// Mark this folder as visited immediately so we never process it again (cycle guard).
	if visited[folderID] && depth > 0 {
		logrus.Debugf("[track-gdrive] cycle detected, skipping already visited folderID=%q", folderID)
		return file.NewDriveDirNode(virtualPath, folderID), nil
	}
	visited[folderID] = true

	entries, err := g.svc.ListFolder(ctx, folderID)
	if err != nil {
		logrus.Debugf("[track-gdrive] ListFolder folderID=%q error: %v", folderID, err)
		return nil, err
	}
	logrus.Debugf("[track-gdrive] depth=%d folderID=%q listed %d entries", depth, folderID, len(entries))

	// Root of this subtree: the folder we're listing
	rootNode := file.NewDriveDirNode(virtualPath, folderID)
	for _, e := range entries {
		childVirtual := path.Join(virtualPath, e.Name)
		if e.IsFolder {
			// Do not recurse into shortcuts (they point to other folders and cause cycles).
			isShortcut := e.MimeType == driveShortcutMimeType
			childNode := file.NewDriveDirNode(childVirtual, e.Id)
			if recursive && !isShortcut && !visited[e.Id] {
				logrus.Debugf("[track-gdrive] recursing into folder %q id=%q", e.Name, e.Id)
				sub, err := g.trackGDriveFolder(ctx, e.Id, childVirtual, true, depth+1, visited)
				if err != nil {
					return nil, err
				}
				childNode.Children = sub.Children
			} else if isShortcut {
				logrus.Debugf("[track-gdrive] skipping shortcut %q id=%q", e.Name, e.Id)
			} else if visited[e.Id] {
				logrus.Debugf("[track-gdrive] skipping already visited folder %q id=%q", e.Name, e.Id)
			}
			rootNode.Children = append(rootNode.Children, childNode)
		} else {
			rootNode.Children = append(rootNode.Children, file.NewDriveFileNode(childVirtual, e.Id))
		}
	}
	return rootNode, nil
}

// NormalizeTrackPath returns the path to use for tracking (strips trailing * and normalizes).
func (g *GDriveScanner) NormalizeTrackPath(pathExp string) (path string, recursive bool) {
	pathExp = strings.TrimSpace(pathExp)
	if strings.HasPrefix(pathExp, file.GDrivePathPrefix) {
		pathExp = strings.TrimPrefix(pathExp, file.GDrivePathPrefix)
	}
	recursive = strings.HasSuffix(pathExp, "*")
	if recursive {
		pathExp = strings.TrimSuffix(pathExp, "*")
		pathExp = strings.TrimSuffix(pathExp, "/")
	}
	pathExp = strings.Trim(pathExp, "/")
	if pathExp == "" {
		pathExp = "/"
	} else {
		pathExp = "/" + pathExp
	}
	return pathExp, recursive
}

// NormalizeGDriveTrackPath is a convenience function that creates a scanner and normalizes the path.
// This is kept for backward compatibility with existing callers.
func NormalizeGDriveTrackPath(pathExp string) (path string, recursive bool) {
	scanner := &GDriveScanner{}
	return scanner.NormalizeTrackPath(pathExp)
}

// GDriveScannableNode is a scannable node for Google Drive.
// TODO: Implement this
type GDriveScannableNode struct {
	strAbsPath string
	cTreeNode  *ds.TreeNode
	children   []ScannableNode
}

// NewGDriveScannableNode creates a new Google Drive scannable node.
func NewGDriveScannableNode(absPath string) *GDriveScannableNode {
	return &GDriveScannableNode{
		strAbsPath: absPath,
	}
}

// EvalNode evaluates the node (currently not implemented).
func (g *GDriveScannableNode) EvalNode(cxt ScannableCxt) error {
	return nil
}

// ConstructTreeNode returns the constructed tree node.
func (g *GDriveScannableNode) ConstructTreeNode() (*ds.TreeNode, error) {
	return g.cTreeNode, nil
}

// GetChildren returns the children nodes and their contexts.
func (g *GDriveScannableNode) GetChildren() ([]ScannableNode, []ScannableCxt, error) {
	return g.children, nil, nil
}
