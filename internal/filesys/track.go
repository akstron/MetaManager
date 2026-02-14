package filesys

import (
	"path/filepath"

	"github.com/heroku/self/MetaManager/internal/cmderror"
	"github.com/heroku/self/MetaManager/internal/ds"
	"github.com/heroku/self/MetaManager/internal/file"
	contextrepo "github.com/heroku/self/MetaManager/internal/repository/filesys"
)

type Tracker interface {
	// Tracks a path and returns a tree node
	Track(path string) (*ds.TreeNode, error)
}

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

	if path[len(path)-1] == '*' {
		return c.trackRec(path)
	}

	return c.trackAbs(path)
}

func (c *ContextAwareTracker) createScannerFromCurrentContext() (Scanner, error) {
	context, err := c.cxtRepo.GetContext()
	if err != nil {
		return nil, err
	}
	contextType, err := c.cxtRepo.GetContextType(context)
	if err != nil {
		return nil, err
	}
	return CreateScannerFromContextType(contextType)
}

func (c *ContextAwareTracker) trackAbs(path string) (*ds.TreeNode, error) {
	// Make sure that this extracted path is tracked
	return file.CreateTreeNodeFromPath(path)
}

func (c *ContextAwareTracker) trackRec(path string) (*ds.TreeNode, error) {
	rootDirPath := path[0 : len(path)-1]
	rootDirPath, err := filepath.Abs(rootDirPath)
	if err != nil {
		return nil, err
	}
	scanner, err := c.createScannerFromCurrentContext()
	if err != nil {
		return nil, err
	}
	return scanner.Scan(rootDirPath)
}
