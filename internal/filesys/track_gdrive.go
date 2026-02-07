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

// TrackGDrive builds a subtree from a Google Drive path. Path is like "/" for root or "/Folder/SubFolder".
// If recursive is true, recurses into all subfolders (up to maxTrackDepth).
// Shortcuts or shared folders that create cycles are skipped to avoid infinite recursion.
func TrackGDrive(ctx context.Context, drivePath string, recursive bool, svc *services.GDriveService) (*ds.TreeNode, error) {
	logrus.Debugf("[track-gdrive] TrackGDrive start path=%q recursive=%v", drivePath, recursive)
	if svc == nil {
		return nil, &cmderror.InvalidOperation{}
	}
	drivePath = strings.Trim(drivePath, "/")
	folderID, err := svc.ResolvePath(ctx, drivePath)
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
	return trackGDriveFolder(ctx, svc, folderID, baseVirtual, recursive, 0, visited)
}

// Drive shortcut MIME type: do not recurse into these (they point to other folders and can create cycles).
const driveShortcutMimeType = "application/vnd.google-apps.shortcut"

func trackGDriveFolder(ctx context.Context, svc *services.GDriveService, folderID, virtualPath string, recursive bool, depth int, visited map[string]bool) (*ds.TreeNode, error) {
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

	entries, err := svc.ListFolder(ctx, folderID)
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
				sub, err := trackGDriveFolder(ctx, svc, e.Id, childVirtual, true, depth+1, visited)
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

// NormalizeGDriveTrackPath returns the path to use for tracking (strips trailing * and normalizes).
func NormalizeGDriveTrackPath(pathExp string) (path string, recursive bool) {
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
