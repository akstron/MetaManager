package filesys

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/heroku/self/MetaManager/internal/file"
	contextrepo "github.com/heroku/self/MetaManager/internal/repository/filesys"
)

type Resolver interface {
	Resolve(path string) (string, error)
}

type BasicResolver struct {
	ctxRepo contextrepo.ContextRepository
}

func NewBasicResolver(ctxRepo contextrepo.ContextRepository) *BasicResolver {
	return &BasicResolver{ctxRepo: ctxRepo}
}

func (r *BasicResolver) Resolve(path string) (string, error) {
	contextName, err := r.ctxRepo.GetContext()
	if err != nil {
		return "", err
	}
	contextType, err := r.ctxRepo.GetContextType(contextName)
	if err != nil {
		return "", err
	}
	switch contextType {
	case contextrepo.TypeGDrive:
		return r.resolveGDrive(path)
	case contextrepo.TypeLocal:
		return r.resolveLocal(path)
	default:
		return "", fmt.Errorf("invalid context type: %s", contextType)
	}
}

func (r *BasicResolver) resolveGDrive(path string) (string, error) {
	cwd, err := r.ctxRepo.GetGDriveCwd()
	if err != nil {
		return "", err
	}
	path = ResolvePath(cwd, path)
	if path == "/" {
		return file.GDrivePathPrefix, nil
	}
	return file.GDrivePathRoot + path, nil
}

func (r *BasicResolver) resolveLocal(path string) (string, error) {
	return filepath.Abs(path)
}

// NormalizePath returns a path like "/" or "/Folder/Sub" (leading slash, no trailing).
func NormalizePath(p string) string {
	p = strings.TrimSpace(p)
	if p == "" || p == "." {
		return "/"
	}
	p = strings.TrimPrefix(p, "/")
	p = strings.TrimSuffix(p, "/")
	if p == "" {
		return "/"
	}
	return "/" + p
}

// ResolvePath resolves target against cwd (absolute path or relative). Returns an absolute path like "/" or "/Folder/Sub".
// If target starts with "/", it is normalized and returned. Otherwise path.Join(cwd, target) is cleaned and normalized.
func ResolvePath(cwd, target string) string {
	if cwd == "" {
		cwd = "/"
	}
	target = strings.TrimSpace(target)
	if target == "" || target == "." {
		return NormalizePath(cwd)
	}
	if strings.HasPrefix(target, "/") {
		return NormalizePath(target)
	}
	// Relative: join with cwd and clean (handles ".." and ".")
	joined := path.Join(cwd, target)
	return NormalizePath(joined)
}
