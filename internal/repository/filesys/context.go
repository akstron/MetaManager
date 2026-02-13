package filesys

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/heroku/self/MetaManager/internal/utils"
)

const (
	contextFileName      = "context"
	contextsJSONFileName = "contexts.json"
	gdriveCwdFileName    = "gdrive_cwd"

	// TypeLocal is the context type for local storage.
	TypeLocal = "local"
	// TypeGDrive is the context type for Google Drive.
	TypeGDrive = "gdrive"
	// ContextEnvVar is the environment variable for the current context. When set, it overrides the context file.
	ContextEnvVar = "MM_CONTEXT"
	// GDriveCwdEnvVar is the environment variable for the current Google Drive directory. When set, overrides the gdrive_cwd file.
	GDriveCwdEnvVar = "MM_GDRIVE_CWD"
)

var validContextTypes = map[string]bool{TypeLocal: true, TypeGDrive: true}

// ContextEntry is a single named context with a type (local or gdrive).
type ContextEntry struct {
	Name string `json:"name"`
	Type string `json:"type"` // "local" or "gdrive"
}

// ContextsFile is the structure stored in contexts.json. Exported for tests that write the file directly.
type ContextsFile struct {
	Contexts []ContextEntry `json:"contexts"`
}

// ContextRepository defines the interface for context storage. Use it for dependency injection and mocking.
type ContextRepository interface {
	LoadContexts() ([]ContextEntry, error)
	AddContext(name, contextType string) error
	SetCurrent(name string) error
	Create(name, contextType string) error
	Delete(name string) error
	DeleteAll() error
	GetContext() (string, error)
	GetContextType(name string) (string, error)
	GetGDriveCwd() (string, error)
	SetGDriveCwd(path string) error
}

// Ensure *ContextRepositoryImpl implements ContextRepository.
var _ ContextRepository = (*ContextRepositoryImpl)(nil)

// BaseDirFunc returns the directory for context files (e.g. next to executable). Inject a mock in tests.
type BaseDirFunc func() (string, error)

// ContextRepositoryImpl holds context state and performs read/write. Use NewContextRepositoryImpl with an optional BaseDirFunc to allow mocking.
type ContextRepositoryImpl struct {
	baseDir BaseDirFunc
}

// NewContextRepositoryImpl creates a ContextRepositoryImpl. If baseDir is nil, the default is used (executable dir, or MM_TEST_CONTEXT_DIR in tests).
func NewContextRepositoryImpl(baseDir BaseDirFunc) *ContextRepositoryImpl {
	if baseDir == nil {
		baseDir = defaultBaseDir
	}
	return &ContextRepositoryImpl{baseDir: baseDir}
}

func defaultBaseDir() (string, error) {
	if d := os.Getenv("MM_TEST_CONTEXT_DIR"); d != "" {
		return d, nil
	}
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("get executable path: %w", err)
	}
	return filepath.Dir(execPath), nil
}

// ContextFilePath returns the path to the context file.
func (s *ContextRepositoryImpl) ContextFilePath() (string, error) {
	dir, err := s.baseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, contextFileName), nil
}

// ContextsJSONPath returns the path to contexts.json.
func (s *ContextRepositoryImpl) ContextsJSONPath() (string, error) {
	dir, err := s.baseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, contextsJSONFileName), nil
}

// gdriveCwdPath returns the path to the gdrive_cwd file.
func (s *ContextRepositoryImpl) gdriveCwdPath() (string, error) {
	dir, err := s.baseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, gdriveCwdFileName), nil
}

// LoadContexts reads and parses contexts.json. Returns a nil slice if the file does not exist.
func (s *ContextRepositoryImpl) LoadContexts() ([]ContextEntry, error) {
	path, err := s.ContextsJSONPath()
	if err != nil {
		return nil, err
	}
	var f ContextsFile
	if err := utils.ReadJSON(path, &f); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return f.Contexts, nil
}

func (s *ContextRepositoryImpl) contextExists(entries []ContextEntry, name string) bool {
	for _, e := range entries {
		if e.Name == name {
			return true
		}
	}
	return false
}

// AddContext adds a context (name and type) to contexts.json. Caller must ensure name is unique.
func (s *ContextRepositoryImpl) AddContext(name, contextType string) error {
	list, err := s.LoadContexts()
	if err != nil {
		return err
	}
	if list == nil {
		list = []ContextEntry{}
	}
	list = append(list, ContextEntry{Name: name, Type: contextType})
	sort.Slice(list, func(i, j int) bool { return list[i].Name < list[j].Name })
	path, err := s.ContextsJSONPath()
	if err != nil {
		return err
	}
	return utils.WriteJSON(path, ContextsFile{Contexts: list}, true)
}

// SetCurrent persists the current context name to the context file. Name must already exist in contexts.json.
func (s *ContextRepositoryImpl) SetCurrent(name string) error {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return fmt.Errorf("context name cannot be empty")
	}
	entries, err := s.LoadContexts()
	if err != nil {
		return err
	}
	if !s.contextExists(entries, name) {
		return fmt.Errorf("context %q not found; create it first with 'context create %s --type local|gdrive'", name, name)
	}
	path, err := s.ContextFilePath()
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(name), 0600); err != nil {
		return fmt.Errorf("write context file %q: %w", path, err)
	}
	return nil
}

// Create adds a new named context with the given type. Returns error if name already exists or type is invalid.
func (s *ContextRepositoryImpl) Create(name, contextType string) error {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return fmt.Errorf("context name cannot be empty")
	}
	contextType = strings.ToLower(strings.TrimSpace(contextType))
	if !validContextTypes[contextType] {
		return fmt.Errorf("type must be %q or %q", TypeLocal, TypeGDrive)
	}
	entries, err := s.LoadContexts()
	if err != nil {
		return err
	}
	if s.contextExists(entries, name) {
		return fmt.Errorf("context %q already exists", name)
	}
	return s.AddContext(name, contextType)
}

// Delete removes a context from contexts.json. If it was the current context (from file), the context file is cleared.
func (s *ContextRepositoryImpl) Delete(name string) error {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return fmt.Errorf("context name cannot be empty")
	}
	entries, err := s.LoadContexts()
	if err != nil {
		return err
	}
	if !s.contextExists(entries, name) {
		return fmt.Errorf("context %q not found", name)
	}
	newEntries := make([]ContextEntry, 0, len(entries)-1)
	for _, e := range entries {
		if e.Name != name {
			newEntries = append(newEntries, e)
		}
	}
	path, err := s.ContextsJSONPath()
	if err != nil {
		return err
	}
	if err := utils.WriteJSON(path, ContextsFile{Contexts: newEntries}, true); err != nil {
		return err
	}
	current, err := s.GetContext()
	if err != nil {
		return err
	}
	if current == name {
		ctxPath, err := s.ContextFilePath()
		if err != nil {
			return err
		}
		if err := os.WriteFile(ctxPath, []byte(""), 0600); err != nil {
			return fmt.Errorf("clear current context file: %w", err)
		}
	}
	return nil
}

// DeleteAll removes all contexts from contexts.json and clears the current context file.
func (s *ContextRepositoryImpl) DeleteAll() error {
	path, err := s.ContextsJSONPath()
	if err != nil {
		return err
	}
	if err := utils.WriteJSON(path, ContextsFile{Contexts: []ContextEntry{}}, true); err != nil {
		return err
	}
	ctxPath, err := s.ContextFilePath()
	if err != nil {
		return err
	}
	if err := os.WriteFile(ctxPath, []byte(""), 0600); err != nil {
		return fmt.Errorf("clear current context file: %w", err)
	}
	return nil
}

// GetContext returns the current context name (env var overrides file). Returns ("", nil) if unset.
func (s *ContextRepositoryImpl) GetContext() (string, error) {
	if v := os.Getenv(ContextEnvVar); v != "" {
		return strings.TrimSpace(strings.ToLower(v)), nil
	}
	path, err := s.ContextFilePath()
	if err != nil {
		return "", err
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(strings.ToLower(string(b))), nil
}

// GetContextType returns the type for the given context name, or ("", nil) if not found.
func (s *ContextRepositoryImpl) GetContextType(name string) (string, error) {
	entries, err := s.LoadContexts()
	if err != nil {
		return "", err
	}
	for _, e := range entries {
		if e.Name == name {
			return e.Type, nil
		}
	}
	return "", nil
}

// GetGDriveCwd returns the current Google Drive working directory (for shell-style navigation).
// Env MM_GDRIVE_CWD overrides the gdrive_cwd file. Returns ("", nil) if unset (meaning root "/").
func (s *ContextRepositoryImpl) GetGDriveCwd() (string, error) {
	if v := os.Getenv(GDriveCwdEnvVar); v != "" {
		return normalizeDrivePath(v), nil
	}
	path, err := s.gdriveCwdPath()
	if err != nil {
		return "", err
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return normalizeDrivePath(strings.TrimSpace(string(b))), nil
}

// SetGDriveCwd persists the current Google Drive path to the gdrive_cwd file.
func (s *ContextRepositoryImpl) SetGDriveCwd(path string) error {
	path = normalizeDrivePath(path)
	filePath, err := s.gdriveCwdPath()
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, []byte(path), 0600)
}

// normalizeDrivePath returns a path like "/" or "/Folder/Sub" (leading slash, no trailing).
func normalizeDrivePath(p string) string {
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

// ResolveGDrivePath resolves target against cwd (absolute path or relative). Returns an absolute path like "/" or "/Folder/Sub".
// If target starts with "/", it is normalized and returned. Otherwise path.Join(cwd, target) is cleaned and normalized.
func ResolveGDrivePath(cwd, target string) string {
	if cwd == "" {
		cwd = "/"
	}
	target = strings.TrimSpace(target)
	if target == "" || target == "." {
		return normalizeDrivePath(cwd)
	}
	if strings.HasPrefix(target, "/") {
		return normalizeDrivePath(target)
	}
	// Relative: join with cwd and clean (handles ".." and ".")
	joined := path.Join(cwd, target)
	return normalizeDrivePath(joined)
}
