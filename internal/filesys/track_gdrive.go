package filesys

import (
	"context"
	"path"
	"strings"

	"github.com/heroku/self/MetaManager/internal/cmderror"
	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/file"
	"github.com/heroku/self/MetaManager/internal/services"
)

const maxTrackDepth = 50

// TrackGDrive builds a subtree from a Google Drive path. Path is like "/" for root or "/Folder/SubFolder".
// If recursive is true, recurses into all subfolders (up to maxTrackDepth).
func TrackGDrive(ctx context.Context, drivePath string, recursive bool, svc *services.GDriveService) (*ds.TreeNode, error) {
	if svc == nil {
		return nil, &cmderror.InvalidOperation{}
	}
	drivePath = strings.Trim(drivePath, "/")
	folderID, err := svc.ResolvePath(ctx, drivePath)
	if err != nil {
		return nil, err
	}
	baseVirtual := file.GDrivePathPrefix
	if drivePath != "" {
		baseVirtual = file.GDrivePathPrefix + drivePath
	}
	return trackGDriveFolder(ctx, svc, folderID, baseVirtual, recursive, 0)
}

func trackGDriveFolder(ctx context.Context, svc *services.GDriveService, folderID, virtualPath string, recursive bool, depth int) (*ds.TreeNode, error) {
	if depth > maxTrackDepth {
		return nil, &cmderror.InvalidOperation{}
	}
	entries, err := svc.ListFolder(ctx, folderID)
	if err != nil {
		return nil, err
	}
	// Root of this subtree: the folder we're listing
	rootNode := file.NewDriveDirNode(virtualPath, folderID)
	for _, e := range entries {
		childVirtual := path.Join(virtualPath, e.Name)
		if e.IsFolder {
			childNode := file.NewDriveDirNode(childVirtual, e.Id)
			if recursive {
				sub, err := trackGDriveFolder(ctx, svc, e.Id, childVirtual, true, depth+1)
				if err != nil {
					return nil, err
				}
				childNode.Children = sub.Children
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
